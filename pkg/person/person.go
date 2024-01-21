package person

type Person struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Nation     string `json:"nation,omitempty"`
}
