package sacloud

// propIcon アイコン内包型
type propIcon struct {
	Icon *Icon // アイコン
}

// SetIconByID 指定のアイコンIDを設定
func (p *propIcon) SetIconByID(id int64) {
	p.Icon = &Icon{Resource: NewResource(id)}
}

// SetIcon 指定のアイコンオブジェクトを設定
func (p *propIcon) SetIcon(icon *Icon) {
	p.Icon = icon
}

// ClearIcon アイコンをクリア(空IDを持つアイコンオブジェクトをセット)
func (p *propIcon) ClearIcon() {
	p.Icon = &Icon{Resource: NewResource(EmptyID)}
}
