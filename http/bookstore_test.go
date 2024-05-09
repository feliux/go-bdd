package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/feliux/go-bdd/http/bookstore"

	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
)

/*
// iter 4 problem
var app = App{store: &bookstore.Store{}}
var e error = nil // for seProduceUnError
*/

// iter 4 solution
type TestContext struct {
	App   App
	Error error
}

// iter 1: scenario 1
func elLibroConIdYTituloHaSidoCreado(ctx context.Context, id, title string) (err error) {
	t := ctx.Value("TestContext").(*TestContext)
	req, err := http.NewRequest(http.MethodGet, "/books/"+id, nil)
	if err != nil {
		return
	}
	w := httptest.NewRecorder()
	t.App.getBook(w, req)
	if w.Code != http.StatusOK {
		return errors.New(fmt.Sprintf("got %d, want %d", w.Code, http.StatusOK))
	}
	var book bookstore.Book
	err = json.NewDecoder(w.Body).Decode(&book)
	if err != nil {
		return
	}
	if book.Title != title {
		return errors.New(fmt.Sprintf("got %s, want %s", book.Title, title))
	}
	return
}

// iter 1: scenario 1
func usuarioCreaUnLibroConLaSiguienteInformacion(ctx context.Context, table *godog.Table) (err error) {
	t := ctx.Value("TestContext").(*TestContext)
	assist := assistdog.NewDefault()
	data, err := assist.ParseMap(table)
	if err != nil {
		return
	}
	book := &bookstore.Book{
		Id:     data["id"],
		Title:  data["title"],
		Author: data["author"],
	}
	bookJson, err := json.Marshal(book)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "/books/", bytes.NewBuffer(bookJson))
	if err != nil {
		return
	}
	w := httptest.NewRecorder()
	t.App.addBook(w, req)
	/*
		// Below if is not true cause scenario 2
		if w.Code != http.StatusAccepted {
			return errors.New(fmt.Sprintf("got %d, want %d", w.Code, http.StatusAccepted))
		}
	*/
	return
}

// iter 2: scenario 2
func laLibreriaContieneLosSiguientesLibros(ctx context.Context, table *godog.Table) (err error) {
	t := ctx.Value("TestContext").(*TestContext)
	assist := assistdog.NewDefault()
	books, err := assist.CreateSlice(new(bookstore.Book), table) // map table to struct
	if err != nil {
		return
	}
	for _, book := range books.([]*bookstore.Book) {
		err = t.App.store.AddBook(book)
		if err != nil {
			t.Error = err // for seProduceUnError
		}
	}
	return
}

// iter 2: scenario 2
func seProduceUnError(ctx context.Context, msg string) (err error) {
	/*
		// Below if is not true cause scenario 3
		if e == nil {
			errors.New(fmt.Sprintf("want %s", msg))
		}
	*/
	t := ctx.Value("TestContext").(*TestContext)
	if t.Error != nil && t.Error.Error() != msg { // There is an error but not which I want to receive
		errors.New(fmt.Sprintf("got %s, want %s", t.Error.Error(), msg))
	}
	return
}

// iter 3: will cause fail with 'godog --random'
// cause scenario 1 not produce error but scenario 2 does
func noSeHaProducidoNingnError(ctx context.Context) error {
	t := ctx.Value("TestContext").(*TestContext)
	if t.Error != nil {
		errors.New(fmt.Sprintf("unexpected error %s", t.Error.Error()))
	}
	return nil
}

// iter 4: feature 2 fails with 'godog --random --concurrency=2'
// cause global vars
func elLibroConIdYTituloHaSidoBorrado(arg1, arg2 string) error {
	return godog.ErrPending
}

// iter 4: feature 2 fails with 'godog --random --concurrency=2'
// cause global vars
func usuarioBorraElLibroConElId(arg1 string) error {
	return godog.ErrPending
}

// iter 4: feature 2 fails with 'godog --random --concurrency=2'
// cause global vars
func usuarioBorraUnLibroConId(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	/*
		// iter 3 solution: reset var e before each scenario to prevent of fails
		ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			e = nil
			return ctx, nil
		})
	*/

	// iter 4 solution: save glbal vars inside a struct and send it to all steps.
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testContext := &TestContext{
			App:   App{store: &bookstore.Store{}},
			Error: nil,
		}
		return context.WithValue(ctx, "TestContext", testContext), nil
	})

	ctx.Step(`^El libro con id "([^"]*)" y titulo "([^"]*)" ha sido creado$`, elLibroConIdYTituloHaSidoCreado)
	ctx.Step(`^Usuario crea un libro con la siguiente informacion:$`, usuarioCreaUnLibroConLaSiguienteInformacion)
	ctx.Step(`^La libreria contiene los siguientes libros:$`, laLibreriaContieneLosSiguientesLibros)
	ctx.Step(`^Se produce un error: "([^"]*)"$`, seProduceUnError)
	ctx.Step(`^No se ha producido ningún error$`, noSeHaProducidoNingnError)
	// iter 4
	ctx.Step(`^El libro con id "([^"]*)" y titulo "([^"]*)" ha sido borrado$`, elLibroConIdYTituloHaSidoBorrado)
	ctx.Step(`^Usuario borra el libro con el id "([^"]*)"$`, usuarioBorraElLibroConElId)
	ctx.Step(`^Usuario borra un libro con id "([^"]*)"$`, usuarioBorraUnLibroConId)
}
