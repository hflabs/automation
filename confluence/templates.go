package confluence

const hashcode_pattern = `.*Эта страница сгенерирована автоматически</ac:parameter>.*?<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>.*?<ac:rich-text-body>\n<p>([0-9a-z]{32})</p>\n</ac:rich-text-body>`

const (
	CheckLine = "<hr />"
	AutoLabel = `
<ac:structured-macro ac:name="info" ac:schema-version="1">
<ac:parameter ac:name="title">Автодокументация</ac:parameter>
<ac:rich-text-body>
<p>Эта страница сгенерирована автоматически</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	NoIncludeIn = `
<ac:structured-macro ac:name="no-include">
<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>
<ac:rich-text-body>
`

	NoIncludeOut = `
</ac:rich-text-body>
</ac:structured-macro>
`

	Expand = `
<ac:structured-macro ac:macro-id="ad79b33f-c1bd-4db2-bb35-1991a0bd1b60" ac:name="expand" ac:schema-version="1">
<ac:rich-text-body>
%s
</ac:rich-text-body>
</ac:structured-macro>
`

	AllTable = `
<table>
<tbody>
%s
</tbody>
</table>
`

	Th  = `<th>%s</th>`
	Sup = `<sup>%s</sup>`

	ThWithStyle = `
<th %s>
%s
</th>
`

	Tr        = "<tr>%s</tr>"
	TableLine = `
<tr %s>
%s
</tr>
`

	Td          = `<td>%s</td>`
	TdWithStyle = `
<td %s>
%s
</td>
`
	Tip = `
<ac:structured-macro ac:name="tip">
<ac:rich-text-body>
<p>%s</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	Restartable = `
<ac:structured-macro ac:name="tip">
<ac:rich-text-body>
<p><a href="https://confluence.hflabs.ru/pages/viewpage.action?pageId=178325752">Возобновляемая</a></p>
</ac:rich-text-body>
</ac:structured-macro>
`
	NotRestartable = `
<ac:structured-macro ac:name="warning">
<ac:rich-text-body>
<p><a href="https://confluence.hflabs.ru/pages/viewpage.action?pageId=178325752">Невозобновляемая</a></p>
</ac:rich-text-body>
</ac:structured-macro>
`
	NotConflicting = `
<ac:structured-macro ac:name="tip">
<ac:rich-text-body>
<p>Неконфликтующая</p>
</ac:rich-text-body>
</ac:structured-macro>
`
	Conflicting = `
<ac:structured-macro ac:name="warning">
<ac:rich-text-body>
<p>Конфликтующая</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	Ul     = `<ul>%s</ul>`
	LiCode = `<li><code>%s;</code></li>`
	Li     = `<li>%s</li>`
	Code   = `<code>%s</code>`

	Style   = `style="%s"`
	Top     = `border-top: 3.0px solid grey;`
	Middle  = `border-left: 3.0px solid grey;border-right: 3.0px solid grey;`
	Bottom  = `border-bottom: 3.0px solid grey;`
	Grey    = `color: grey;`
	ColSpan = `colspan="%s"`
	RowSpan = `rowspan="%s"`

	Children = `<ac:structured-macro ac:name="children">
<ac:parameter ac:name="all">true</ac:parameter>
<ac:parameter ac:name="excerptType">simple</ac:parameter>
</ac:structured-macro>`

	ChildrenShort = `<ac:structured-macro ac:name="children">
<ac:parameter ac:name="excerptType">simple</ac:parameter>
</ac:structured-macro>`
	HrefToOutward = `<a href="%s">%s</a>`
	Href          = `<a href="https://confluence.hflabs.ru/pages/viewpage.action?pageId=%s">%s</a>`

	HrefByName = `<a href="https://confluence.hflabs.ru/display/%s/%s">%s</a>`

	H2  = `<h2>%s</h2>`
	H3  = `<h3>%s</h3>`
	H4  = `<h4>%s</h4>`
	Par = `<p>%s</p>`

	Yes  = `<ac:emoticon ac:name="tick"></ac:emoticon>`
	Link = `<ac:link>
<ri:page ri:content-title="%s" ri:space-key="%s"></ri:page>
<ac:plain-text-link-body><![CDATA[%s]]></ac:plain-text-link-body>
</ac:link>`

	SpaceLink = `<a href="https://confluence.hflabs.ru/display/%s">%s</a>`

	TdCenter = `<td style="text-align: center">%s</td>`

	Include = `
<p>
<ac:structured-macro ac:name="include">
<ac:parameter ac:name="">
<ac:link>
<ri:page ri:space-key="%s" ri:content-title="%s" />
</ac:link>
</ac:parameter>
</ac:structured-macro>
</p>
`

	Hide = `
<ac:structured-macro ac:name="hfl" ac:schema-version="1">
<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>
<ac:rich-text-body>
%s
</ac:rich-text-body>
</ac:structured-macro>
`
	HideHash = `
<ac:structured-macro ac:name="hfl" ac:schema-version="1">
<ac:parameter ac:name="hiddentitle">Эта страница сгенерирована автоматически</ac:parameter>
<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>
<ac:rich-text-body>
<p>%s</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	BottomLine = `
<p>*Указывает на особенности обработки поля в SOAP.</p>
<p>**Указывает на поведение реквизита в интерфейсе «Единого Клиента».</p>
<ul>
<li>
Полнотекстовый поиск — участвует в полнотекстовом поиске в веб-интерфейсе и SOAP.
</li>
<li>
Расширенный поиск — участвует только в
<a href="https://confluence.hflabs.ru/pages/viewpage.action?pageId=%s">расширенном поиске</a>
в веб-интерфейсе и SOAP.
</li>
</ul>
`

	Excerpt = `
<ac:structured-macro ac:name="excerpt" ac:schema-version="1">
<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>
<ac:rich-text-body>
%s
</ac:rich-text-body>
</ac:structured-macro>
`

	ExcerptHidden = `
<ac:structured-macro ac:name="excerpt" ac:schema-version="1">
<ac:parameter ac:name="hidden">true</ac:parameter>
<ac:parameter ac:name="atlassian-macro-output-type">INLINE</ac:parameter>
<ac:rich-text-body>
%s
</ac:rich-text-body>
</ac:structured-macro>
`

	WarnForDatabaseTable = `
<ac:structured-macro ac:name="warning" ac:schema-version="1">
<ac:rich-text-body>
<p>
<strong>Это описание внутренней таблицы «Единого клиента»,
структура меняется без предупреждения и уведомления</strong>
</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	Info = `
<ac:structured-macro ac:macro-id="680e9bb3-f0a5-445f-bed2-55e16e83e212" ac:name="info" ac:schema-version="1">
<ac:rich-text-body>
<p>%s</p>
</ac:rich-text-body>
</ac:structured-macro>
`

	Toc = `<ac:structured-macro ac:name="toc" ac:schema-version="1"/>`

	RefHead = `<strong> <span style="font-size: 19.0px;"> %s</span> </strong>`
)
