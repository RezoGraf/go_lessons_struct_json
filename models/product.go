package models

//Попробовал одну большую структуру но вернулся
//к раздельной конструкции,
//в связи с тем что так значительно проще оказалось
//выбирать данные из дочерних таблиц БД
//когда структурно данные поделены

// type Product struct {
// 	Title *string `json:"title" db:"title"`
// 	Tags  []struct {
// 		ID  int    `json:"id" db:"id"`
// 		Tag string `json:"tag" db:"tag"`
// 	} `json:"tags" db:"tags"`
// 	Description    string `json:"description,omitempty" db:"description"`
// 	Price          string `json:"price,omitempty" db:"price"`
// 	Additionalinfo struct {
// 		Title   string `json:"title" db:"title"`
// 		Comment string `json:"comment" db:"comment"`
// 	} `json:"additionalinfo,omitempty" db:"additionalinfo"`
// }

type Product struct {
	Title          *string        `json:"title" db:"title"`
	Tags           []Tag          `json:"tags"`
	Description    string         `json:"description,omitempty" db:"description"`
	Price          string         `json:"price,omitempty" db:"price"`
	Additionalinfo Additionalinfo `json:"additionalinfo,omitempty"`
}
type Tag struct {
	ID  int    `json:"id" db:"id"`
	Tag string `json:"tag" db:"tag"`
}
type Additionalinfo struct {
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
}
