package pdf

type Description struct {
	TitleInfo TitleInfo `xml:"title-info"`
}

type TitleInfo struct {
	Author     Author `xml:"author"`
	Title      string `xml:"book-title"`
	Annotation struct {
		Content string `xml:",innerxml"`
	} `xml:"annotation"`
	Date string `xml:"date"`
}

type Author struct {
	FirstName string `xml:"first-name"`
	LastName  string `xml:"last-name"`
}

type Pg struct {
	Text string `xml:",innerxml"`
}

type Title struct {
	Pgs []Pg `xml:"p"`
}

type SubTitle struct {
	Text string `xml:",innerxml"`
}
