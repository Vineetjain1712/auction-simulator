package auction

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vineetjain1712/auction-simulator/internal/models"
)

// ItemGenerator generates random auction items
type ItemGenerator struct {
	rand *rand.Rand
	mu   sync.Mutex // Protects rand for thread-safety
}

// NewItemGenerator creates a new item generator
func NewItemGenerator() *ItemGenerator {
	// Create a new random source with current time as seed
	source := rand.NewSource(time.Now().UnixNano())
	return &ItemGenerator{
		rand: rand.New(source),
	}
}

// Predefined lists for generating realistic items
var (
	categories     = []string{"Electronics", "Art", "Collectibles", "Jewelry", "Furniture", "Books", "Clothing"}
	brands         = []string{"Apple", "Samsung", "Sony", "Nike", "Canon", "Rolex", "Generic"}
	conditions     = []string{"New", "Like New", "Used", "Refurbished", "Fair"}
	colors         = []string{"Black", "White", "Silver", "Gold", "Blue", "Red", "Green"}
	sizes          = []string{"Small", "Medium", "Large", "XL", "XXL", "One Size"}
	materials      = []string{"Metal", "Plastic", "Wood", "Leather", "Fabric", "Glass", "Ceramic"}
	rarities       = []string{"Common", "Uncommon", "Rare", "Very Rare", "Ultra Rare"}
	origins        = []string{"USA", "China", "Japan", "Germany", "Italy", "France", "UK"}
	certifications = []string{"CE", "FCC", "ISO9001", "RoHS", "None", "UL", "Energy Star"}
)

// GenerateItem creates a random auction item with all 20 attributes
func (g *ItemGenerator) GenerateItem(id int) models.AuctionItem {
	g.mu.Lock()
	defer g.mu.Unlock()

	category := g.randomChoiceUnsafe(categories)
	brand := g.randomChoiceUnsafe(brands)

	// Generate a contextual name based on category
	name := fmt.Sprintf("%s %s %d", brand, category, id)

	return models.AuctionItem{
		// Attribute 1-5
		ID:        id,
		Name:      name,
		Category:  category,
		Brand:     brand,
		Condition: g.randomChoiceUnsafe(conditions),

		// Attribute 6-10
		Color:    g.randomChoiceUnsafe(colors),
		Size:     g.randomChoiceUnsafe(sizes),
		Weight:   g.randomFloatUnsafe(0.1, 50.0), // 0.1kg to 50kg
		Material: g.randomChoiceUnsafe(materials),
		YearMade: g.randomIntUnsafe(2010, 2024),

		// Attribute 11-15
		Origin:      g.randomChoiceUnsafe(origins),
		Rarity:      g.randomChoiceUnsafe(rarities),
		BasePrice:   g.randomFloatUnsafe(10.0, 5000.0), // $10 to $5000
		Description: fmt.Sprintf("High quality %s from %s", category, brand),
		Features:    fmt.Sprintf("Premium %s with excellent quality", category),

		// Attribute 16-20
		Warranty:      g.randomIntUnsafe(0, 36), // 0 to 36 months
		ShipWeight:    g.randomFloatUnsafe(0.2, 55.0),
		Dimensions:    fmt.Sprintf("%.1fx%.1fx%.1f", g.randomFloatUnsafe(5, 100), g.randomFloatUnsafe(5, 100), g.randomFloatUnsafe(5, 100)),
		Certification: g.randomChoiceUnsafe(certifications),
		Rating:        g.randomFloatUnsafe(3.0, 10.0), // 3.0 to 10.0
	}
}

// Helper functions - "Unsafe" means caller must hold mutex

func (g *ItemGenerator) randomChoiceUnsafe(choices []string) string {
	return choices[g.rand.Intn(len(choices))]
}

func (g *ItemGenerator) randomIntUnsafe(min, max int) int {
	return min + g.rand.Intn(max-min+1)
}

func (g *ItemGenerator) randomFloatUnsafe(min, max float64) float64 {
	return min + g.rand.Float64()*(max-min)
}

// GenerateItems generates multiple items at once
func (g *ItemGenerator) GenerateItems(count int) []models.AuctionItem {
	items := make([]models.AuctionItem, count)
	for i := range count {
		items[i] = g.GenerateItem(i + 1)
	}
	return items
}
