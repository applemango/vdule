# About
Vtuberの配信予定やVtuberの情報を取得するAPIです

> [!IMPORTANT]
> 現時点ではHololive所属のVtuberは自動で取得 ( 15分更新 ) されますが、それ以外のVtuberは取得されません
> 
> セクション`Add Vtuber`を参考にVtuberを登録してください、30分に一回更新されます

> [!TIP]
> Youtube Data Apiのレートリミットがデフォルトでは月10000回なので、大体6人分登録可能です
> 5分に一回更新に変更し1人分に変えることも可能です

# How to use
1. `.env`ファイルを`.env.example`を参考に作る
2. dbのtableを作る

# Example

## Example Response
### /schedule
```json
{
    "videos": [
        {
            "channel": {
                "banner_image": "https://yt3.googleusercontent.com/ibCEuiR9Za7MNWSVsGXPdQEL2ZVWPXwLYA1nTlWtxf0X_0-vZlKqK_OBS_hJkjRWAXzvSB8u5Qg",
                "description": "こんこよ〜！ホロライブ所属、秘密結社holoXの頭脳🧪博衣こよりだよ〜(ᐡ •͈ ·̫ •͈ ᐡ)ﾉｼ💕\n初見古参関係なく、みんなで楽しめるチャンネルを作っていきたいです！\n「夢を全部叶えること」が夢！\nチャンネル登録、Twitterのフォローよろしくお願いします！✨\n\nKonkoyo~!! My name is Hakui Koyori!!\nI'm Hololive member, The brains of the secret society holoX!!\nI will create a channel that everyone can enjoy!\nSubscribe to my channel and follow me on Twitter!\n\nママ：ももこ先生 @momoco_haru\nパパ：rariemonn先生 @rariemonn765\n\n＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊\n\n　🐶🧪ハッシュタグ🧪🐶\n　イラスト：#こよりすけっち\n　生放送関連：#こより実験中\n　ファンネーム：こよりの助手くん\n\n＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊\n\n　🧪メンバーシップはじまりました❣🧪\n　コメントでスタンプを使えたり、メンバーシップ限定の配信や動画、コミュニティ投稿、壁紙配布があります！\n　▼登録はここから▼\n　https://www.youtube.com/channel/UC6eWCld0KwmyHFbAqK3V-Rw/join\n\n＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊\n\n　お手紙はこちら💌\n\n　〒173-0003\n　東京都板橋区加賀1丁目6番1号　\n　ネットデポ新板橋\n　カバー株式会社 ホロライブ プレゼント係分\n　博衣こより宛\n\n　※お約束ごと→ https://www.hololive.tv/contact ※\n\n＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊ … ⬡ …＊ … ⬡ … ＊\n",
                "handle": "hakuikoyori",
                "icon_image": "https://yt3.ggpht.com/WO7ItKNmy6tW_NQ82g8c1y74CZSw6GsSdynsE5s2csuEok2fHRrAaGcBV3JJO-2BxEOXXA8lvw=s800-c-k-c0x00ffffff-no-rj",
                "id": "UC6eWCld0KwmyHFbAqK3V-Rw",
                "name": "Koyori ch. 博衣こより - holoX -",
                "publish_at": "2021-09-12T03:04:25.578989Z",
                "subscriber_count": 1070000,
                "trailer": "",
                "uploads_playlist": "UU6eWCld0KwmyHFbAqK3V-Rw",
                "video_count": 1529,
                "view_count": 296003258
            },
            "date": {
                "day": 11,
                "hour": 11,
                "month": 2,
                "year": 2024
            },
            "id": "fKtoJIX6eks",
            "thumbnail": "https://img.youtube.com/vi/fKtoJIX6eks/mqdefault.jpg",
            "title": "【龍が如く8】9章なう！メインもサブも良すぎて何回泣くんだ今作😭 #6 【博衣こより/ホロライブ】【ネタバレあり】"
        },
        ...
    ]
}
```

## Schema
```sql
DROP TABLE IF EXISTS vtuber_organization;
CREATE TABLE vtuber_organization (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE,
    description TEXT
);
DROP TABLE IF EXISTS vtuber_tag;
CREATE TABLE vtuber_tag (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);
DROP TABLE IF EXISTS vtuber;
CREATE TABLE vtuber (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id VARCHAR(16),
    handle TEXT UNIQUE,
    name TEXT,
    description TEXT,
    organization_id INTEGER,
    is_crawl BOOLEAN
);
DROP TABLE IF EXISTS vtubers_tag;
CREATE TABLE vtubers_tag (
    vtuber_id INTEGER,
    tag_id INTEGER
);
DROP TABLE IF EXISTS channel;
CREATE TABLE channel (
    id VARCHAR(16) PRIMARY KEY,
    handle TEXT,
    view_count INTEGER,
    subscriber_count INTEGER,
    video_count INTEGER
);
DROP TABLE IF EXISTS schedule;
CREATE TABLE schedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    thumbnail TEXT,
    channel_id VARCHAR(16),
    handle TEXT,
    title TEXT,
    unix_time INTEGER,
    time_year INTEGER,
    time_month INTEGER,
    time_day INTEGER,
    time_hour INTEGER
);
DROP TABLE IF EXISTS video;
CREATE TABLE video (
    id VARCHAR(16) PRIMARY KEY,
    channel_id VARCHAR(16),
    handle TEXT,
    title TEXT,
    thumbnail TEXT,
    is_live BOOLEAN,
    is_now_on_air BOOLEAN,
    live_scheduled_start_year INTEGER,
    live_scheduled_start_month INTEGER,
    live_scheduled_start_day INTEGER,
    live_scheduled_start_hour INTEGER,
    live_scheduled_start_minute INTEGER
);
```

## Add Vtuber
```go
vtuber.RegisterOrganization(vtuber.RegisterOrganizationProps{
	Name:        "にゃんたじあ！",
	Description: "",
})
err := vtuber.RegisterVtuber(vtuber.RegisterVtuberProps{
	Handle:         "@nyamafujianzu",
	Tags:           []string{"女性"},
	Name:           "若魔藤あんず",
	Description:    "",
	OrganizationId: 1,
	IsCrawl:        true,
})
```