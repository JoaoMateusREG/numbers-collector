package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

type Registro struct {
	CPF    string `json:"cpf"`
	Numero int64  `json:"numero"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./dados.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Tabela
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS registros (
		cpf TEXT PRIMARY KEY,
		numero INTEGER
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	http.HandleFunc("/registro", setupCORS(manipularRegistro))

	fmt.Println("Servidor rodando na porta 7531...")
	log.Fatal(http.ListenAndServe(":7531", nil)) // [3]
}

func manipularRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed) // [2]
		return
	}

	var reg Registro
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest) // [2]
		return
	}

	// 1. Normalizar o CPF: Remover tudo que não for número antes de validar e salvar
    re := regexp.MustCompile(`[^0-9]`)
    reg.CPF = re.ReplaceAllString(reg.CPF, "") // Agora reg.CPF terá apenas números

    // 2. Validar o CPF (agora já limpo)
    if !validarCPF(reg.CPF) {
        http.Error(w, "CPF inválido", http.StatusBadRequest)
        return
    }


	// 1. Valida se o CPF é matematicamente real
	if !validarCPF(reg.CPF) {
		http.Error(w, "CPF inválido", http.StatusBadRequest) // [1], [2]
		return
	}

	// 2. Valida se o número tem exatamente 11 dígitos
	// Convertemos para string para contar o comprimento
	if !validarNumero(reg.Numero) {
		http.Error(w, "O numero deve ter exatamente 11 dígitos", http.StatusBadRequest) // [1], [2]
		return
	}


	query := `
	INSERT INTO registros(cpf, numero) VALUES(?, ?)
	ON CONFLICT(cpf) DO UPDATE SET numero=excluded.numero;
	`
	_, err := db.Exec(query, reg.CPF, reg.Numero)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco", http.StatusInternalServerError) // [2]
		log.Println("Erro SQL:", err)
		return
	}

	w.WriteHeader(http.StatusOK) // [1]
	w.Write([]byte("Registro válido salvo com sucesso"))
}

// validarNumero verifica se o inteiro possui 11 dígitos.
// Aviso: Como o campo é int (ou int64), zeros à esquerda (ex: 012...) são ignorados.
func validarNumero(n int64) bool {
	s := strconv.FormatInt(n, 10)
	return len(s) == 11
}

// validarCPF aplica o algoritmo padrão de verificação de dígitos.
// Nota: Esta lógica é algoritmo padrão (matemático) e não depende das bibliotecas importadas.
func validarCPF(cpf string) bool {
    re := regexp.MustCompile(`[^0-9]`)
    cpf = re.ReplaceAllString(cpf, "")

    if len(cpf) != 11 {
        return false
    }

    // Verifica se todos os dígitos são iguais (ex: 111.111.111-11)
    iguais := true
    for i := 1; i < 11; i++ {
        if cpf[i] != cpf[0] {
            iguais = false
            break
        }
    }
    if iguais {
        return false
    }

    // Cálculo do 1º dígito verificador
    soma := 0
    for i := 0; i < 9; i++ {
        soma += int(cpf[i]-'0') * (10 - i)
    }
    resto := (soma * 10) % 11
    if resto == 10 { resto = 0 }
    
    if resto != int(cpf[9]-'0') { // Comparando com o 10º dígito
        return false
    }

    // Cálculo do 2º dígito verificador
    soma = 0
    for i := 0; i < 10; i++ {
        soma += int(cpf[i]-'0') * (11 - i)
    }
    resto = (soma * 10) % 11
    if resto == 10 { resto = 0 }

    if resto != int(cpf[10]-'0') { // Comparando com o 11º dígito
        return false
    }

    return true
}

func setupCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Permite qualquer origem (em produção, você pode trocar "*" pelo seu domínio)
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
        // Permite os métodos que sua API usa
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        
        // Permite cabeçalhos personalizados como Content-Type
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Importante: Responde imediatamente a requisições "Preflight" (OPTIONS)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Passa a requisição para o seu handler original (manipularRegistro)
        next(w, r)
    }
}