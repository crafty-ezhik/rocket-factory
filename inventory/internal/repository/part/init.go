package part

import (
	"context"
	"log"
	"time"
)

// Init - наполняет хранилище n случайными запчастями
func (r *repository) Init() {
	parts := GetParts()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, part := range parts {
		_, err := r.collection.InsertOne(ctx, part)
		if err != nil {
			log.Printf("Ошибка вставки детали %s: %v\n", part.Name, err)
		} else {
			log.Printf("✅ Добавлена деталь: %s\n", part.Name)
		}
	}
}
