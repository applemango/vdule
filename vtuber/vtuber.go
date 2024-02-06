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
	Id           int32               `json:"id,omitempty"`
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

type Row interface {
	*sql.Row | *sql.Rows
	Scan(dest ...any) error
}

func GetVtuberCoreHelper[T Row](row T) (*Vtuber, error) {
	var (
		vtuber         Vtuber
		organizationId sql.NullInt32
	)
	err := row.Scan(&vtuber.Id, &vtuber.Handle, &vtuber.Name, &vtuber.Description, &organizationId, &vtuber.IsCrawl)
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
	return &vtuber, nil
}

func GetVtuberCore(query string, args ...any) (*Vtuber, error) {
	row := db.Conn.QueryRow(query, args...)
	return GetVtuberCoreHelper(row)
}

func GetVtuberByHandle(handle string) (*Vtuber, error) {
	return GetVtuberCore(`SELECT id, handle, name, description, organization_id, is_crawl FROM vtuber WHERE handle = ?`, youtube.ParseChannelHandle(handle))
}

func GetVtubersCore(query string, args ...any) ([]Vtuber, error) {
	var vs []Vtuber
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		v, err := GetVtuberCoreHelper(rows)
		if err != nil {
			return nil, err
		}
		vs = append(vs, *v)
	}
	return vs, nil
}

func GetVtubers() ([]Vtuber, error) {
	return GetVtubersCore(`SELECT id, handle, name, description, organization_id, is_crawl FROM vtuber`)
}

type RegisterVtuberProps struct {
	Handle         string   `json:"handle,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	OrganizationId int32    `json:"organization_id,omitempty"`
	IsCrawl        bool     `json:"is_crawl,omitempty"`
}

func GetTagByName(name string) (*Tag, error) {
	var tag Tag
	row := db.Conn.QueryRow(`SELECT id, name FROM vtuber_tag WHERE name = ?`, name)
	err := row.Scan(&tag.Id, &tag.Name)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func CreateTag(name string) (*Tag, error) {
	_, err := db.Conn.Exec(`INSERT INTO vtuber_tag (name) VALUES (?)`, name)
	if err != nil {
		return nil, err
	}
	return GetTagByName(name)
}

func GetTagByNameForce(name string) (*Tag, error) {
	if tag, err := GetTagByName(name); err == nil {
		return tag, nil
	}
	return CreateTag(name)
}

func GetTagsByNamesForce(names []string) ([]Tag, error) {
	var tags []Tag
	for _, name := range names {
		tag, err := GetTagByNameForce(name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func (v Vtuber) AppendTag(tag Tag) error {
	_, err := db.Conn.Exec(`INSERT INTO vtubers_tag (vtuber_id, tag_id) VALUES (?, ?)`, v.Id, tag.Id)
	return err
}

func (v Vtuber) ResetTag(tags []Tag) error {
	_, err := db.Conn.Exec(`DELETE FROM vtubers_tag WHERE vtuber_id = ?`, v.Id)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		if err := v.AppendTag(tag); err != nil {
			return err
		}
	}
	return nil
}

func RegisterVtuber(props RegisterVtuberProps) error {
	channel, err := youtube.T.GetChannelByHandle(props.Handle)
	if err != nil {
		return err
	}
	if props.OrganizationId != -1 {
		_, err = db.Conn.Exec(`INSERT INTO vtuber (channel_id, handle, name, description, organization_id, is_crawl) VALUES (?, ?, ?, ?, ?, ?)`, channel.Id, youtube.ParseChannelHandle(channel.Handle), props.Name, props.Description, props.OrganizationId, props.IsCrawl)
	} else {
		_, err = db.Conn.Exec(`INSERT INTO vtuber (channel_id, handle, name, description, is_crawl) VALUES (?, ?, ?, ?, ?)`, channel.Id, youtube.ParseChannelHandle(channel.Handle), props.Name, props.Description, props.IsCrawl)
	}
	if err != nil {
		return err
	}
	tags, err := GetTagsByNamesForce(props.Tags)
	if err != nil {
		return err
	}
	vtuber, err := GetVtuberByHandle(channel.Handle)
	if err != nil {
		return err
	}
	err = vtuber.ResetTag(tags)
	if err != nil {
		return err
	}
	return nil
}

type RegisterOrganizationProps struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func RegisterOrganization(props RegisterOrganizationProps) error {
	_, err := db.Conn.Exec(`INSERT INTO vtuber_organization (name, description) VALUES (?, ?)`, props.Name, props.Description)
	return err
}
