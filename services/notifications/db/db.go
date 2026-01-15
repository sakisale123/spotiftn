package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func Init() {
	cassandraHost := os.Getenv("CASSANDRA_HOST")
	if cassandraHost == "" {
		cassandraHost = "localhost"
	}

	var cluster *gocql.ClusterConfig
	var err error

	for i := 0; i < 30; i++ {
		cluster = gocql.NewCluster(cassandraHost)
		cluster.Keyspace = "system"
		cluster.Consistency = gocql.Quorum
		cluster.ProtoVersion = 4
		cluster.Timeout = 10 * time.Second

		Session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		log.Printf("Waiting for Cassandra... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to Cassandra:", err)
	}

	fmt.Println("Connected to Cassandra")

	createKeyspace()

	Session.Close()
	cluster.Keyspace = "notifications_db"
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Could not connect to keyspace:", err)
	}

	createTable()
	seedData()
}

func createKeyspace() {
	query := `CREATE KEYSPACE IF NOT EXISTS notifications_db 
	WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`

	if err := Session.Query(query).Exec(); err != nil {
		log.Fatal("Failed to create keyspace:", err)
	}
	fmt.Println("Keyspace 'notifications_db' ensured")
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS notifications (
		id UUID,
		user_id text,
		type text,
		message text,
		created_at timestamp,
		is_read boolean,
		PRIMARY KEY (user_id, created_at)
	) WITH CLUSTERING ORDER BY (created_at DESC);`

	if err := Session.Query(query).Exec(); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table 'notifications' ensured")
}

func seedData() {
	var count int

	err := Session.Query("SELECT count(*) FROM notifications WHERE user_id = ?", "global").Scan(&count)
	if err == nil && count > 0 {
		fmt.Println("Global data already exists, skipping seed.")
		return
	}

	fmt.Println("Seeding global data...")

	titles := []string{
		"New Release: The Weeknd - After Hours",
		"Concert Alert: Arctic Monkeys in your area",
		"Trending: 'Levitating' by Dua Lipa",
		"Genre Spotlight: Techno Minimal",
		"Recommended: Tame Impala - The Slow Rush",
		"Weekly Top 10: Check out what's hot",
		"New Album: 'Midnights' by Taylor Swift",
		"Artist to Watch: Fred again..",
		"Festival Season: Glastonbury Tickets Available",
		"Flashback Friday: The Beatles - Abbey Road",
	}

	types := []string{"release", "concert", "trending", "recommendation"}

	for i, title := range titles {
		id, _ := gocql.RandomUUID()
		notifType := types[i%len(types)]

		createdAt := time.Now().Add(-time.Duration(i) * time.Hour)

		err := Session.Query(`INSERT INTO notifications (id, user_id, type, message, created_at, is_read) VALUES (?, ?, ?, ?, ?, ?)`,
			id, "global", notifType, title, createdAt, false).Exec()
		if err != nil {
			log.Printf("Failed to seed notification: %v", err)
		}
	}
	fmt.Println("Seeding completed.")
}
