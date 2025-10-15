package repository

type InventoryRepository interface {
	Get()
	List()
	Init()
}
