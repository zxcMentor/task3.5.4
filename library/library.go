package library

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golibrary/utils"

	"github.com/Masterminds/squirrel"
)

type LibraryFacade struct {
	DB *sql.DB
}

func NewLibraryFacade(db *sql.DB) *LibraryFacade {
	return &LibraryFacade{DB: db}
}

func (lf *LibraryFacade) StartLibrary() {
	userCount, bookCount, authorCount := lf.getCounts()
	lf.generateAndInsertDataIfNeeded(userCount, bookCount, authorCount)
}

func (lf *LibraryFacade) PrintLibraryUsers() {
	users, err := lf.getUsers()
	if err != nil {
		log.Println("Ошибка при получении пользователей:", err)
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Name)
		fmt.Println("Арендованные книги:")
		for _, book := range user.RentedBooks {
			fmt.Printf("  ID: %d, Название: %s\n", book.ID, book.Name)
		}
		fmt.Println("---------------")
	}
}

func (lf *LibraryFacade) getCounts() (int, int, int) {
	counts := make([]int, 3)
	tables := []string{"users", "books", "authors"}

	for i, table := range tables {
		query := squirrel.Select("COUNT(*)").From(table)
		sql, args, err := query.ToSql()
		if err != nil {
			log.Fatal(err)
		}

		err = lf.DB.QueryRow(sql, args...).Scan(&counts[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	return counts[0], counts[1], counts[2]
}

func (lf *LibraryFacade) generateAndInsertDataIfNeeded(userCount, bookCount, authorCount int) (string, error) {
	var result strings.Builder

	if userCount == 0 {
		users := utils.GenerateAndInsertUsers(lf.DB, 50, bookCount)
		result.WriteString(fmt.Sprintf("Сгенерировано и добавлено пользователей: %d\n", len(users)))
	}

	if bookCount == 0 {
		books := utils.GenerateAndInsertBooks(lf.DB, 100)
		result.WriteString(fmt.Sprintf("Сгенерировано и добавлено книг: %d\n", len(books)))
	}

	if authorCount == 0 {
		authors := utils.GenerateAndInsertAuthors(lf.DB, 10, bookCount)
		result.WriteString(fmt.Sprintf("Сгенерировано и добавлено авторов: %d\n", len(authors)))
	}

	return result.String(), nil
}

func (lf *LibraryFacade) getBookIDs() ([]int, error) {
	query := squirrel.Select("id").From("books")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := lf.DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookIDs []int
	for rows.Next() {
		var bookID int
		err := rows.Scan(&bookID)
		if err != nil {
			return nil, err
		}
		bookIDs = append(bookIDs, bookID)
	}

	return bookIDs, nil
}

func (lf *LibraryFacade) getUsers(bookIDs []int) ([]utils.User, error) {
	query := squirrel.Select("id", "name").From("users")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := lf.DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []utils.User
	for rows.Next() {
		var user utils.User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}

		user.RentedBooks = utils.GetRandomBooks(lf.DB, bookIDs)

		users = append(users, user)
	}

	return users, nil
}
