# chyoa
Change your own adventure written in Golang













WEB VERSION

Роут GET api/cyoa/{chapter}, 
function GetChapterHandler(w, *r)
{
Взять chapter из запроса
Вызвать getChapter
Валидировать
В шаблон загрузить значения chapter.Title, chapter.Text
для каждого из chapter.Options загрузить в шаблон <a href=(api/cyoa/{chapter.options.arc)>chapter.options.Text</a>
Если chapter.Options пуст, то другой шаблон с предложением "Thans for playing"
}
