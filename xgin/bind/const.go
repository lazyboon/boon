package bind

type Key string

const (
	KeyJSON          Key = "_lazyboon.bind.json.key"
	KeyXML           Key = "_lazyboon.bind.xml.key"
	KeyForm          Key = "_lazyboon.bind.form.key"
	KeyQuery         Key = "_lazyboon.bind.query.key"
	KeyFormPost      Key = "_lazyboon.bind.form_post.key"
	KeyFormMultipart Key = "_lazyboon.bind.form_multipart.key"
	KeyProtoBuf      Key = "_lazyboon.bind.proto_buf.key"
	KeyMsgPack       Key = "_lazyboon.bind.msg_pack.key"
	KeyYAML          Key = "_lazyboon.bind.yaml.key"
	KeyHeader        Key = "_lazyboon.bind.header.key"
	KeyTOML          Key = "_lazyboon.bind.toml.key"
	KeyUri           Key = "_lazyboon.bind.uri.key"
)
