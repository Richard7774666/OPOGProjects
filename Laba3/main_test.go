package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGymTracker(t *testing.T) {
	t.Run("AddVisit and CountVisitsByClient", func(t *testing.T) {
		tracker := NewGymTracker()

		// Success flow тести
		t.Run("success flow - count visits for single client", func(t *testing.T) {
			clientName := "Іван Петренко"
			visit := Visit{
				ID:             1,
				Client:         clientName,
				MembershipType: "Стандарт",
				Date:           time.Now(),
			}

			tracker.AddVisit(visit)
			tracker.AddVisit(visit) // Додаємо ще одне відвідування

			count := tracker.CountVisitsByClient(clientName)
			require.Equal(t, 2, count, "Кількість відвідувань має бути 2")
		})

		// Edge case тести
		t.Run("edge case - count visits for non-existent client", func(t *testing.T) {
			count := tracker.CountVisitsByClient("Неіснуючий Клієнт")
			require.Equal(t, 0, count, "Для неіснуючого клієнта має повертатися 0")
		})
	})

	t.Run("GetMostFrequentClient", func(t *testing.T) {
		tracker := NewGymTracker()

		// Підготовка тестових даних
		clients := []struct {
			name   string
			visits int
			mType  string
		}{
			{"Клієнт А", 3, "Стандарт"},
			{"Клієнт Б", 5, "Преміум"},
			{"Клієнт В", 1, "Базовий"},
		}

		for _, c := range clients {
			for i := 0; i < c.visits; i++ {
				tracker.AddVisit(Visit{
					ID:             i + 1,
					Client:         c.name,
					MembershipType: c.mType,
					Date:           time.Now().AddDate(0, 0, -i),
				})
			}
		}

		// Success flow тест
		t.Run("success flow - find most frequent client", func(t *testing.T) {
			name, visits := tracker.GetMostFrequentClient()
			require.Equal(t, "Клієнт Б", name, "Найчастіший клієнт має бути 'Клієнт Б'")
			require.Equal(t, 5, visits, "Кількість відвідувань має бути 5")
		})
	})
}

// Table-driven тести для CountVisitsByClient
func TestCountVisitsByClient_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		clientName    string
		addVisits     []Visit
		expectedCount int
	}{
		{
			name:          "no visits",
			clientName:    "Клієнт X",
			addVisits:     []Visit{},
			expectedCount: 0,
		},
		{
			name:       "single visit",
			clientName: "Клієнт Y",
			addVisits: []Visit{
				{ID: 1, Client: "Клієнт Y", MembershipType: "Стандарт", Date: time.Now()},
			},
			expectedCount: 1,
		},
		{
			name:       "multiple visits",
			clientName: "Клієнт Z",
			addVisits: []Visit{
				{ID: 1, Client: "Клієнт Z", MembershipType: "Стандарт", Date: time.Now()},
				{ID: 2, Client: "Клієнт Z", MembershipType: "Преміум", Date: time.Now()},
				{ID: 3, Client: "Інший Клієнт", MembershipType: "Базовий", Date: time.Now()},
			},
			expectedCount: 2,
		},
		{
			name:       "visits with different membership types",
			clientName: "Клієнт W",
			addVisits: []Visit{
				{ID: 1, Client: "Клієнт W", MembershipType: "Стандарт", Date: time.Now()},
				{ID: 2, Client: "Клієнт W", MembershipType: "Преміум", Date: time.Now()},
				{ID: 3, Client: "Клієнт W", MembershipType: "Базовий", Date: time.Now()},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewGymTracker()
			for _, visit := range tt.addVisits {
				tracker.AddVisit(visit)
			}

			count := tracker.CountVisitsByClient(tt.clientName)
			require.Equal(t, tt.expectedCount, count,
				"Неправильна кількість відвідувань для клієнта %s", tt.clientName)
		})
	}
}
