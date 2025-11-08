package main

import (
	"fmt"
	"log"

	"github.com/iqmalr-pedia/go-product/internal/models"
	"github.com/iqmalr-pedia/go-product/pkg/database"
	"gorm.io/gorm"
)

func main() {
	database.ConnectDB()
	db := database.GetDB()

	// Declare variables outside the transaction to access them later
	var (
		mainCategories              []models.Category
		electronicsSubCategories    []models.Category
		computersSubCategories      []models.Category
		smartphonesSubCategories    []models.Category
		fashionSubCategories        []models.Category
		homeGardenSubCategories     []models.Category
		sportsOutdoorsSubCategories []models.Category
		healthBeautySubCategories   []models.Category
		toysGamesSubCategories      []models.Category
		booksMediaSubCategories     []models.Category
		automotiveSubCategories     []models.Category
		foodBeveragesSubCategories  []models.Category
		officeSuppliesSubCategories []models.Category
	)

	err := database.WithTransaction(db, func(tx *gorm.DB) error {
		// Clear existing categories
		if err := tx.Exec("DELETE FROM categories").Error; err != nil {
			return fmt.Errorf("failed to clear existing categories: %w", err)
		}
		log.Println("Cleared existing categories")

		// Create main categories
		mainCategories = []models.Category{
			{
				Name:      "Electronics",
				Desc:      "Electronic devices, gadgets, and accessories",
				Icon:      "electronics",
				ImageURL:  "https://example.com/images/electronics.jpg",
				Level:     1,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				Name:      "Fashion",
				Desc:      "Clothing, shoes, and accessories for men, women, and children",
				Icon:      "fashion",
				ImageURL:  "https://example.com/images/fashion.jpg",
				Level:     1,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				Name:      "Home & Garden",
				Desc:      "Furniture, decor, and garden supplies",
				Icon:      "home",
				ImageURL:  "https://example.com/images/home-garden.jpg",
				Level:     1,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				Name:      "Sports & Outdoors",
				Desc:      "Sports equipment, activewear, and outdoor gear",
				Icon:      "sports",
				ImageURL:  "https://example.com/images/sports-outdoors.jpg",
				Level:     1,
				SortOrder: 4,
				IsActive:  true,
			},
			{
				Name:      "Health & Beauty",
				Desc:      "Personal care, cosmetics, and health products",
				Icon:      "health",
				ImageURL:  "https://example.com/images/health-beauty.jpg",
				Level:     1,
				SortOrder: 5,
				IsActive:  true,
			},
			{
				Name:      "Toys & Games",
				Desc:      "Toys, games, and hobbies for all ages",
				Icon:      "toys",
				ImageURL:  "https://example.com/images/toys-games.jpg",
				Level:     1,
				SortOrder: 6,
				IsActive:  true,
			},
			{
				Name:      "Books & Media",
				Desc:      "Books, movies, music, and other media",
				Icon:      "books",
				ImageURL:  "https://example.com/images/books-media.jpg",
				Level:     1,
				SortOrder: 7,
				IsActive:  true,
			},
			{
				Name:      "Automotive",
				Desc:      "Car parts, accessories, and tools",
				Icon:      "automotive",
				ImageURL:  "https://example.com/images/automotive.jpg",
				Level:     1,
				SortOrder: 8,
				IsActive:  true,
			},
			{
				Name:      "Food & Beverages",
				Desc:      "Groceries, specialty foods, and beverages",
				Icon:      "food",
				ImageURL:  "https://example.com/images/food-beverages.jpg",
				Level:     1,
				SortOrder: 9,
				IsActive:  true,
			},
			{
				Name:      "Office Supplies",
				Desc:      "Stationery, office furniture, and equipment",
				Icon:      "office",
				ImageURL:  "https://example.com/images/office-supplies.jpg",
				Level:     1,
				SortOrder: 10,
				IsActive:  true,
			},
		}

		// Save main categories
		for i := range mainCategories {
			if err := tx.Create(&mainCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create main category %s: %w", mainCategories[i].Name, err)
			}
		}

		// Electronics subcategories
		electronicsSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "Computers",
				Desc:      "Laptops, desktops, and computer accessories",
				Icon:      "computers",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "Smartphones",
				Desc:      "Mobile phones and accessories",
				Icon:      "smartphones",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "TV & Home Theater",
				Desc:      "Televisions, sound systems, and home entertainment",
				Icon:      "tv",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "Cameras",
				Desc:      "Digital cameras, lenses, and photography equipment",
				Icon:      "cameras",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "Audio",
				Desc:      "Headphones, speakers, and audio equipment",
				Icon:      "audio",
				Level:     2,
				SortOrder: 5,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[0].ID,
				Name:      "Wearables",
				Desc:      "Smartwatches, fitness trackers, and other wearable tech",
				Icon:      "wearables",
				Level:     2,
				SortOrder: 6,
				IsActive:  true,
			},
		}

		// Save Electronics subcategories
		var computersID, smartphonesID uint
		for i := range electronicsSubCategories {
			if err := tx.Create(&electronicsSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create electronics subcategory %s: %w", electronicsSubCategories[i].Name, err)
			}
			if electronicsSubCategories[i].Name == "Computers" {
				computersID = electronicsSubCategories[i].ID
			}
			if electronicsSubCategories[i].Name == "Smartphones" {
				smartphonesID = electronicsSubCategories[i].ID
			}
		}

		// Computers sub-subcategories
		computersSubCategories = []models.Category{
			{
				ParentID:  &computersID,
				Name:      "Laptops",
				Desc:      "Notebook computers and laptops",
				Icon:      "laptops",
				Level:     3,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &computersID,
				Name:      "Desktops",
				Desc:      "Desktop computers and all-in-ones",
				Icon:      "desktops",
				Level:     3,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &computersID,
				Name:      "Tablets",
				Desc:      "Tablet computers and e-readers",
				Icon:      "tablets",
				Level:     3,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &computersID,
				Name:      "Computer Accessories",
				Desc:      "Keyboards, mice, monitors, and other accessories",
				Icon:      "computer-accessories",
				Level:     3,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Computers sub-subcategories
		for i := range computersSubCategories {
			if err := tx.Create(&computersSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create computers subcategory %s: %w", computersSubCategories[i].Name, err)
			}
		}

		// Smartphones sub-subcategories
		smartphonesSubCategories = []models.Category{
			{
				ParentID:  &smartphonesID,
				Name:      "Android Phones",
				Desc:      "Smartphones running Android OS",
				Icon:      "android",
				Level:     3,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &smartphonesID,
				Name:      "iPhones",
				Desc:      "Apple iPhone smartphones",
				Icon:      "iphone",
				Level:     3,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &smartphonesID,
				Name:      "Phone Accessories",
				Desc:      "Cases, chargers, screen protectors, and other accessories",
				Icon:      "phone-accessories",
				Level:     3,
				SortOrder: 3,
				IsActive:  true,
			},
		}

		// Save Smartphones sub-subcategories
		for i := range smartphonesSubCategories {
			if err := tx.Create(&smartphonesSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create smartphones subcategory %s: %w", smartphonesSubCategories[i].Name, err)
			}
		}

		// Fashion subcategories
		fashionSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Men's Clothing",
				Desc:      "Shirts, pants, jackets, and other clothing for men",
				Icon:      "mens-clothing",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Women's Clothing",
				Desc:      "Dresses, tops, bottoms, and other clothing for women",
				Icon:      "womens-clothing",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Kids' Clothing",
				Desc:      "Clothing for boys and girls",
				Icon:      "kids-clothing",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Shoes",
				Desc:      "Footwear for men, women, and children",
				Icon:      "shoes",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Bags & Accessories",
				Desc:      "Handbags, wallets, belts, and other accessories",
				Icon:      "bags-accessories",
				Level:     2,
				SortOrder: 5,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[1].ID,
				Name:      "Jewelry",
				Desc:      "Necklaces, earrings, rings, and other jewelry",
				Icon:      "jewelry",
				Level:     2,
				SortOrder: 6,
				IsActive:  true,
			},
		}

		// Save Fashion subcategories
		for i := range fashionSubCategories {
			if err := tx.Create(&fashionSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create fashion subcategory %s: %w", fashionSubCategories[i].Name, err)
			}
		}

		// Home & Garden subcategories
		homeGardenSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[2].ID,
				Name:      "Furniture",
				Desc:      "Indoor and outdoor furniture",
				Icon:      "furniture",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[2].ID,
				Name:      "Home Decor",
				Desc:      "Decorative items and accents for the home",
				Icon:      "home-decor",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[2].ID,
				Name:      "Kitchen & Dining",
				Desc:      "Cookware, dinnerware, and kitchen appliances",
				Icon:      "kitchen-dining",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[2].ID,
				Name:      "Bedding & Bath",
				Desc:      "Linens, towels, and bathroom accessories",
				Icon:      "bedding-bath",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[2].ID,
				Name:      "Garden & Outdoor",
				Desc:      "Plants, tools, and outdoor living supplies",
				Icon:      "garden-outdoor",
				Level:     2,
				SortOrder: 5,
				IsActive:  true,
			},
		}

		// Save Home & Garden subcategories
		for i := range homeGardenSubCategories {
			if err := tx.Create(&homeGardenSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create home & garden subcategory %s: %w", homeGardenSubCategories[i].Name, err)
			}
		}

		// Sports & Outdoors subcategories
		sportsOutdoorsSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[3].ID,
				Name:      "Exercise & Fitness",
				Desc:      "Fitness equipment, apparel, and accessories",
				Icon:      "fitness",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[3].ID,
				Name:      "Team Sports",
				Desc:      "Equipment for team sports like soccer, basketball, etc.",
				Icon:      "team-sports",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[3].ID,
				Name:      "Outdoor Recreation",
				Desc:      "Camping, hiking, and outdoor adventure gear",
				Icon:      "outdoor-recreation",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[3].ID,
				Name:      "Water Sports",
				Desc:      "Swimming, diving, and water sports equipment",
				Icon:      "water-sports",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Sports & Outdoors subcategories
		for i := range sportsOutdoorsSubCategories {
			if err := tx.Create(&sportsOutdoorsSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create sports & outdoors subcategory %s: %w", sportsOutdoorsSubCategories[i].Name, err)
			}
		}

		// Health & Beauty subcategories
		healthBeautySubCategories = []models.Category{
			{
				ParentID:  &mainCategories[4].ID,
				Name:      "Skin Care",
				Desc:      "Facial cleansers, moisturizers, and treatments",
				Icon:      "skin-care",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[4].ID,
				Name:      "Makeup",
				Desc:      "Cosmetics and beauty products",
				Icon:      "makeup",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[4].ID,
				Name:      "Hair Care",
				Desc:      "Shampoos, conditioners, and styling products",
				Icon:      "hair-care",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[4].ID,
				Name:      "Personal Care",
				Desc:      "Bath products, deodorants, and toiletries",
				Icon:      "personal-care",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[4].ID,
				Name:      "Health Care",
				Desc:      "Vitamins, supplements, and health monitoring devices",
				Icon:      "health-care",
				Level:     2,
				SortOrder: 5,
				IsActive:  true,
			},
		}

		// Save Health & Beauty subcategories
		for i := range healthBeautySubCategories {
			if err := tx.Create(&healthBeautySubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create health & beauty subcategory %s: %w", healthBeautySubCategories[i].Name, err)
			}
		}

		// Toys & Games subcategories
		toysGamesSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[5].ID,
				Name:      "Toys",
				Desc:      "Action figures, dolls, and other toys",
				Icon:      "toys",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[5].ID,
				Name:      "Games & Puzzles",
				Desc:      "Board games, card games, and puzzles",
				Icon:      "games-puzzles",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[5].ID,
				Name:      "Learning & Education",
				Desc:      "Educational toys and learning materials",
				Icon:      "learning-education",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[5].ID,
				Name:      "Video Games",
				Desc:      "Video games, consoles, and accessories",
				Icon:      "video-games",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Toys & Games subcategories
		for i := range toysGamesSubCategories {
			if err := tx.Create(&toysGamesSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create toys & games subcategory %s: %w", toysGamesSubCategories[i].Name, err)
			}
		}

		// Books & Media subcategories
		booksMediaSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[6].ID,
				Name:      "Books",
				Desc:      "Fiction, non-fiction, and educational books",
				Icon:      "books",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[6].ID,
				Name:      "Movies & TV",
				Desc:      "DVDs, Blu-rays, and digital movies",
				Icon:      "movies-tv",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[6].ID,
				Name:      "Music",
				Desc:      "CDs, vinyl records, and digital music",
				Icon:      "music",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[6].ID,
				Name:      "Magazines",
				Desc:      "Print and digital magazines",
				Icon:      "magazines",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Books & Media subcategories
		for i := range booksMediaSubCategories {
			if err := tx.Create(&booksMediaSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create books & media subcategory %s: %w", booksMediaSubCategories[i].Name, err)
			}
		}

		// Automotive subcategories
		automotiveSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[7].ID,
				Name:      "Parts & Accessories",
				Desc:      "Car parts, fluids, and accessories",
				Icon:      "parts-accessories",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[7].ID,
				Name:      "Tools & Equipment",
				Desc:      "Automotive tools and diagnostic equipment",
				Icon:      "tools-equipment",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[7].ID,
				Name:      "Car Care",
				Desc:      "Cleaning supplies, waxes, and car care products",
				Icon:      "car-care",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[7].ID,
				Name:      "Car Electronics", // FIXED: Changed from "Electronics" to "Car Electronics"
				Desc:      "Car audio, GPS, and other electronics",
				Icon:      "car-electronics",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Automotive subcategories
		for i := range automotiveSubCategories {
			if err := tx.Create(&automotiveSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create automotive subcategory %s: %w", automotiveSubCategories[i].Name, err)
			}
		}

		// Food & Beverages subcategories
		foodBeveragesSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[8].ID,
				Name:      "Groceries",
				Desc:      "Everyday food items and pantry staples",
				Icon:      "groceries",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[8].ID,
				Name:      "Beverages",
				Desc:      "Coffee, tea, soft drinks, and alcoholic beverages",
				Icon:      "beverages",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[8].ID,
				Name:      "Snacks & Sweets",
				Desc:      "Chips, candy, and other snacks",
				Icon:      "snacks-sweets",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[8].ID,
				Name:      "Specialty Foods",
				Desc:      "Gourmet, organic, and specialty food items",
				Icon:      "specialty-foods",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Food & Beverages subcategories
		for i := range foodBeveragesSubCategories {
			if err := tx.Create(&foodBeveragesSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create food & beverages subcategory %s: %w", foodBeveragesSubCategories[i].Name, err)
			}
		}

		// Office Supplies subcategories
		officeSuppliesSubCategories = []models.Category{
			{
				ParentID:  &mainCategories[9].ID,
				Name:      "Stationery",
				Desc:      "Paper, pens, and writing supplies",
				Icon:      "stationery",
				Level:     2,
				SortOrder: 1,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[9].ID,
				Name:      "Office Furniture",
				Desc:      "Desks, chairs, and other office furniture",
				Icon:      "office-furniture",
				Level:     2,
				SortOrder: 2,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[9].ID,
				Name:      "Office Electronics",
				Desc:      "Printers, scanners, and other office electronics",
				Icon:      "office-electronics",
				Level:     2,
				SortOrder: 3,
				IsActive:  true,
			},
			{
				ParentID:  &mainCategories[9].ID,
				Name:      "Storage & Organization",
				Desc:      "Filing cabinets, boxes, and organization solutions",
				Icon:      "storage-organization",
				Level:     2,
				SortOrder: 4,
				IsActive:  true,
			},
		}

		// Save Office Supplies subcategories
		for i := range officeSuppliesSubCategories {
			if err := tx.Create(&officeSuppliesSubCategories[i]).Error; err != nil {
				return fmt.Errorf("failed to create office supplies subcategory %s: %w", officeSuppliesSubCategories[i].Name, err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}

	log.Println("Category seeding completed successfully!")
	fmt.Println("Created categories:")
	fmt.Printf("- %d main categories\n", len(mainCategories))
	fmt.Printf("- %d Electronics subcategories\n", len(electronicsSubCategories))
	fmt.Printf("- %d Computers sub-subcategories\n", len(computersSubCategories))
	fmt.Printf("- %d Smartphones sub-subcategories\n", len(smartphonesSubCategories))
	fmt.Printf("- %d Fashion subcategories\n", len(fashionSubCategories))
	fmt.Printf("- %d Home & Garden subcategories\n", len(homeGardenSubCategories))
	fmt.Printf("- %d Sports & Outdoors subcategories\n", len(sportsOutdoorsSubCategories))
	fmt.Printf("- %d Health & Beauty subcategories\n", len(healthBeautySubCategories))
	fmt.Printf("- %d Toys & Games subcategories\n", len(toysGamesSubCategories))
	fmt.Printf("- %d Books & Media subcategories\n", len(booksMediaSubCategories))
	fmt.Printf("- %d Automotive subcategories\n", len(automotiveSubCategories))
	fmt.Printf("- %d Food & Beverages subcategories\n", len(foodBeveragesSubCategories))
	fmt.Printf("- %d Office Supplies subcategories\n", len(officeSuppliesSubCategories))
}
