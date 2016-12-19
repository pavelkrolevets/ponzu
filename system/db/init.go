package db

import (
	"encoding/json"
	"log"

	"github.com/ponzu-cms/ponzu/system/admin/config"
	"github.com/ponzu-cms/ponzu/system/item"

	"github.com/boltdb/bolt"
	"github.com/nilslice/jwt"
)

var store *bolt.DB

// Close exports the abillity to close our db file. Should be called with defer
// after call to Init() from the same place.
func Close() {
	err := store.Close()
	if err != nil {
		log.Println(err)
	}
}

// Init creates a db connection, initializes db with required info, sets secrets
func Init() {
	var err error
	store, err = bolt.Open("system.db", 0666, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = store.Update(func(tx *bolt.Tx) error {
		// initialize db with all content type buckets & sorted bucket for type
		for t := range item.Types {
			_, err := tx.CreateBucketIfNotExists([]byte(t))
			if err != nil {
				return err
			}

			_, err = tx.CreateBucketIfNotExists([]byte(t + "__sorted"))
			if err != nil {
				return err
			}
		}

		// init db with other buckets as needed
		buckets := []string{"__config", "__users", "__contentIndex"}
		for _, name := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
		}

		// seed db with configs structure if not present
		b := tx.Bucket([]byte("__config"))
		if b.Get([]byte("settings")) == nil {
			j, err := json.Marshal(&config.Config{})
			if err != nil {
				return err
			}

			err = b.Put([]byte("settings"), j)
			if err != nil {
				return err
			}
		}

		clientSecret := ConfigCache("client_secret")

		if clientSecret != "" {
			jwt.Secret([]byte(clientSecret))
		}

		return nil
	})
	if err != nil {
		log.Fatalln("Coudn't initialize db with buckets.", err)
	}

	// invalidate cache on system start
	err = InvalidateCache()
	if err != nil {
		log.Fatalln("Failed to invalidate cache.", err)
	}

	go func() {
		for t := range item.Types {
			SortContent(t)
		}
	}()

}

// SystemInitComplete checks if there is at least 1 admin user in the db which
// would indicate that the system has been configured to the minimum required.
func SystemInitComplete() bool {
	complete := false

	err := store.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("__users"))

		err := users.ForEach(func(k, v []byte) error {
			complete = true
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		complete = false
		log.Fatalln(err)
	}

	return complete
}
