package main

import (
	"fmt"
	"time"
)

// Membership представляє абонемент у спортзал
type Membership struct {
	ID       int
	Type     string
	Duration int // у місяцях
}

// Visit представляє відвідування спортзалу
type Visit struct {
	ID             int
	Client         string
	MembershipType string
	Date           time.Time
}

// GymTracker відстежує відвідування спортзалу
type GymTracker struct {
	visitsByMembership map[string][]Visit
	allVisits          []Visit
}

// NewGymTracker створює новий трекер відвідувань
func NewGymTracker() *GymTracker {
	return &GymTracker{
		visitsByMembership: make(map[string][]Visit),
	}
}

// AddVisit додає нове відвідування до трекера
func (gt *GymTracker) AddVisit(v Visit) {
	gt.allVisits = append(gt.allVisits, v)
	gt.visitsByMembership[v.MembershipType] = append(gt.visitsByMembership[v.MembershipType], v)
}

// CountVisitsByClient повертає кількість відвідувань для конкретного клієнта
func (gt *GymTracker) CountVisitsByClient(clientName string) int {
	count := 0
	for _, visit := range gt.allVisits {
		if visit.Client == clientName {
			count++
		}
	}
	return count
}

// GetMostFrequentClient повертає клієнта з найбільшою кількістю відвідувань
func (gt *GymTracker) GetMostFrequentClient() (string, int) {
	clientVisits := make(map[string]int)

	for _, visit := range gt.allVisits {
		clientVisits[visit.Client]++
	}

	var topClient string
	maxVisits := 0

	for client, visits := range clientVisits {
		if visits > maxVisits {
			maxVisits = visits
			topClient = client
		}
	}

	return topClient, maxVisits
}

func main() {
	// Ініціалізація трекера
	tracker := NewGymTracker()

	// Додавання тестових даних
	tracker.AddVisit(Visit{
		ID:             1,
		Client:         "Іван Петренко",
		MembershipType: "Стандарт",
		Date:           time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
	})

	tracker.AddVisit(Visit{
		ID:             2,
		Client:         "Марія Коваль",
		MembershipType: "Преміум",
		Date:           time.Date(2023, time.January, 16, 0, 0, 0, 0, time.UTC),
	})

	tracker.AddVisit(Visit{
		ID:             3,
		Client:         "Іван Петренко",
		MembershipType: "Стандарт",
		Date:           time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
	})

	tracker.AddVisit(Visit{
		ID:             4,
		Client:         "Олексій Сідоров",
		MembershipType: "Базовий",
		Date:           time.Date(2023, time.February, 5, 0, 0, 0, 0, time.UTC),
	})

	tracker.AddVisit(Visit{
		ID:             5,
		Client:         "Марія Коваль",
		MembershipType: "Преміум",
		Date:           time.Date(2023, time.March, 10, 0, 0, 0, 0, time.UTC),
	})

	// Тестування функціоналу
	fmt.Println("Кількість відвідувань Івана Петренка:", tracker.CountVisitsByClient("Іван Петренко"))
	fmt.Println("Кількість відвідувань Марії Коваль:", tracker.CountVisitsByClient("Марія Коваль"))

	topClient, visits := tracker.GetMostFrequentClient()
	fmt.Printf("Найчастіший відвідувач: %s (%d відвідувань)\n", topClient, visits)

	// Виведення відвідувань за типами абонементів
	fmt.Println("\nВідвідування за типами абонементів:")
	for mType, visits := range tracker.visitsByMembership {
		fmt.Printf("- %s: %d відвідувань\n", mType, len(visits))
	}
}
