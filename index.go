package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	pgx "github.com/jackc/pgx/v4"
)

func connectDb() (*pgx.Conn, error) { //detecta a URL de conexão ao banco de dados
	connstr := os.Getenv("DATATESTE_URL")
	if len(connstr) == 0 {
		err := fmt.Errorf("sem url de conexão")
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), connstr) //conecta à URL do banco de dados
	if err != nil {
		fmt.Errorf("impossivel estaelecer conexão[%v]: %v", connstr, err)
	}
	return conn, nil
}

type Amigo struct {
	Cod       uuid.UUID `json: cod`
	Nome      string    `json: nome`
	Sobrenome string    `json: sobrenome`
	Telefone  string    `json: telefone`
	Cidade    string    `json: cidade`
}

type JsonResponse struct {
	Type    string  `json: "type"`
	Data    []Amigo `json "data"`
	Message string  `json: "message"`
}

var conn *pgx.Conn

func init() {
	var err error
	conn, err = connectDb()
	if err != nil {
		panic(fmt.Sprintf("FUDEU: %v", err))
	}

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/amigos", getAllFriends).Methods("GET")
	router.HandleFunc("/amigos/{nomeAmigo}", getAFriend).Methods("GET")

	port := "8080"
	//SERV THE APP
	fmt.Printf("Server at %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getAFriend(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("inicio GetAFriend\n")
	vars := mux.Vars(r)
	var sql string = "SELECT cod, nome, sobrenome, telefone, cidade FROM contato WHERE nome = $1"
	var user Amigo
	nomeAmigo := vars["nomeAmigo"]
	err := conn.QueryRow(context.Background(), sql, nomeAmigo).Scan(
		&user.Cod,
		&user.Nome,
		&user.Sobrenome,
		&user.Telefone,
		&user.Cidade,
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("GetAFriend: cp2\n")

	//amigos := []Amigo{user}
	bytes, err := json.Marshal(user)
	if err != nil {
		w.Write([]byte("falhou ao gerar o json"))
	}
	fmt.Printf("GetAFriend: cp 3\n")
	w.Write(bytes)
	fmt.Printf("Deve ter enviado %v", user)
}

func getAllFriends(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("inicio GetAFriend\n")
	var sql string = "SELECT cod, nome, sobrenome, telefone, cidade FROM contato"
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var amigos []Amigo

	for rows.Next() {
		var user Amigo
		rows.Scan(
			&user.Cod,
			&user.Nome,
			&user.Sobrenome,
			&user.Telefone,
			&user.Cidade,
		)
		amigos = append(amigos, user)
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("GetAFriend: cp2\n")

	//amigos := []Amigo{user}
	bytes, err := json.Marshal(amigos)
	if err != nil {
		w.Write([]byte("falhou ao gerar o json"))
	}
	fmt.Printf("GetAFriend: cp 3\n")
	w.Write(bytes)
	fmt.Printf("Deve ter enviado %v", amigos)
}
