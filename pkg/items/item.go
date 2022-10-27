package items

type Item struct {
	ID          uint32 `schema:"-"`
	Title       string `schema:"title,required"`
	Description string `schema:"description,required"`
	CreatedBy   uint32 `schema:"-"`
}

type ItemsRepo interface {
	GetAll() ([]*Item, error)
	GetByID(id uint32) (*Item, error)
	Add(item *Item) (uint32, error)
	Update(newItem *Item) (bool, error)
	Delete(id uint32) (bool, error)
}
