package models

type Post struct {
	ID		 int
	Year     string `json:"year"`
	Type     string `json:"type"`
	Country  string `json:"country"`
	Activity string `json:"activity"`
	Age      string `json:"age"`
	Fatal    string `json:"fatal_y_n"`
}
