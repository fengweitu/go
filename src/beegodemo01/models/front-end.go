package models

type RegisterUser struct {
	Username string
	Password string
}

type LoginUser struct {
	Username string
	Password string
	Check    bool
}

type ArticleContnet struct {
	ArticleTitle   string
	ArticleType    string
	ArticleContent string
	ImgName        string
	ImgContent     string
}

type Paging struct {
	Handle      string
	ArticleType string
	NowPage     string
}

type DelArticle struct {
	Id          string
	NowPage     string
	ArticleType string
}

type DelType struct {
	Id string
}
type Detail struct {
	Id       string
	TypeName string
}
type Edit struct {
	Id         string
	Title      string
	Content    string
	ImgName    string
	ImgContent string
}
