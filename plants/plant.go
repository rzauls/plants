package plants

type Plant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
}

func (p Plant) Valid() map[string]string {
	problems := make(map[string]string)
	if p.Name == "" {
		problems["name"] = "name cannot be empty"
	}

	if p.Height < 0 {
		problems["height"] = "height cannot be negative"
	}

	return problems
}
