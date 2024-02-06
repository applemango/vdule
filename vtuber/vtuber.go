package vtuber

import (
	"database/sql"
	db "vdule/utils/db/sqlite3"
	"vdule/vtuber/youtube"
)

type Tag struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Organization struct {
	Id          int32  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Vtuber struct {
	Id           string              `json:"id,omitempty"`
	Channel      youtube.TubeChannel `json:"channel"`
	Handle       string              `json:"handle,omitempty"`
	Tags         []Tag               `json:"tags,omitempty"`
	Name         string              `json:"name,omitempty"`
	Description  string              `json:"description,omitempty"`
	Organization *Organization       `json:"organization,omitempty"`
	IsCrawl      bool                `json:"is_crawl,omitempty"`
}

func (v Vtuber) GetVtubersTag() ([]Tag, error) {
	var tags []Tag
	tagsRow, err := db.Conn.Query(`SELECT vt.id, vt.name FROM vtubers_tag JOIN main.vtuber_tag vt on vtubers_tag.tag_id = vt.id WHERE vtuber_id = ?`, v.Id)
	if err != nil {
		return nil, err
	}
	for tagsRow.Next() {
		var tag Tag
		err = tagsRow.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func GetVtubersOrganization(id int32) (*Organization, error) {
	var organization Organization
	row := db.Conn.QueryRow(`SELECT id, name, description FROM vtuber_organization WHERE id = ?`, id)
	err := row.Scan(&organization.Id, &organization.Name, &organization.Description)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func GetVtubers() ([]Vtuber, error) {
	rows, err := db.Conn.Query(`SELECT id, handle, name, description, organization_id, is_crawl FROM vtuber`)
	if err != nil {
		return nil, err
	}
	var vtubers []Vtuber
	for rows.Next() {
		var (
			vtuber         Vtuber
			organizationId sql.NullInt32
		)
		err = rows.Scan(&vtuber.Id, &vtuber.Handle, &vtuber.Name, &vtuber.Description, &organizationId, &vtuber.IsCrawl)
		if err != nil {
			return nil, err
		}
		if organizationId.Valid {
			if organization, err := GetVtubersOrganization(organizationId.Int32); err == nil {
				vtuber.Organization = organization
			}
		}
		c, _ := youtube.GetRawChannelByHandleCache(vtuber.Handle)
		vtuber.Channel = youtube.ChannelToTubeChannel(c)
		if tags, err := vtuber.GetVtubersTag(); err == nil {
			vtuber.Tags = tags
		}
		vtubers = append(vtubers, vtuber)
	}
	return vtubers, nil
}
