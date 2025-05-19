# Pagination Package (Go)

Este projeto implementa um sistema genérico de paginação para consultas SQL no Go, com suporte para paginação via parâmetros e retorno em formato JSON estruturado.

## 📦 Estrutura do Projeto

- `main.go`: Exemplo de uso real com acesso a banco de dados PostgreSQL.
- `abstract.go`: Lógica base da paginação genérica.
- `interface.go`: Interface `IPagination` que define os métodos necessários.
- `envelope.go`: Definição do tipo `Envelope`, um `map[string]interface{}` para retornar resultados.
- `validation.go`: Regras de validação para a paginação.
- `*_test.go`: Testes unitários com `sqlmock`.

---

## 🚀 Como Usar

### 1. Estrutura esperada

Implemente um tipo de struct representando seu modelo de dados:

```go
type TabelaExame struct {
    IDModulo int    `json:"id_modulo" db:"id_modulo"`
    Modulo   string `json:"modulo" db:"modulo"`
    IDExame  int    `json:"id_exame" db:"id_exame"`
    Exame    string `json:"exame" db:"exame"`
}
```

### 2. Inicialize o banco

```go
var banco DB

func init() {
    banco.Init("postgresql://usuario:senha@localhost:5430/nome_do_banco")
}
```

### 3. Execute a paginação

```go
pag := pagination.New[TabelaExame](10, 4, banco.Conn)

query := `
    SELECT modulo.id_modulo, modulo, id_exame, exame
    FROM ep_dw.exame
    JOIN ep_dw.modulo ON modulo.id_modulo = exame.id_modulo
    WHERE exame = $1
`
params := []interface{}{"Glicose"}

pag.SetOrder("id_modulo")
pag.SetRawQuery(query, params...)

data, err := pag.JSON(func(e *[]TabelaExame, rows *sql.Rows) error {
    return banco.ScanRowsToStruct(rows, e)
})

fmt.Println(data, err)
```

---

## 📄 Retorno

A função `JSON(...)` retorna um `Envelope` com a seguinte estrutura:

```json
{
  "data": [ ... ],
  "next_page": 5,
  "previous_page": 3,
  "count": 40
}
```

---

## ✅ Testes

Execute os testes com:

```bash
go test ./...
```

Testes unitários utilizam `sqlmock` para simular interações com o banco.

---

## 📌 Requisitos

- Go 1.18+ (usa generics)
- PostgreSQL (ou outro banco com driver `database/sql`)
- Biblioteca para mock: [`github.com/DATA-DOG/go-sqlmock`](https://github.com/DATA-DOG/go-sqlmock)
- [`testify`](https://github.com/stretchr/testify) para asserts nos testes

---

## 🛠️ Extensibilidade

Você pode reaproveitar o sistema com qualquer struct, bastando especificar o tipo ao chamar `pagination.New[T]()`. O método `SetRawQuery(...)` permite usar qualquer query parametrizada.