# 🚀 HareID API - Guia de Instalação e Execução

Este documento descreve o processo passo a passo para configurar, instalar as dependências e rodar a API HareID localmente.

## 1. Pré-requisitos

Certifique-se de ter as seguintes ferramentas instaladas em sua máquina:

*   **Go (Golang)**: Versão 1.24.0 ou superior.
    *   Verifique com: `go version`
*   **Git**: Para clonar o repositório.
*   **Editor de Código**: Recomendado VS Code (com a extensão Go instalada).

## 2. Configuração do Ambiente (.env)

A aplicação utiliza um arquivo `.env` para carregar variáveis sensíveis e configurações de porta.

1.  Na **raiz do projeto**, crie um arquivo chamado `.env` baseando-se no arquivo de exemplo `.env.example`.
2.  O conteúdo deve seguir este modelo:

```env
SUPABASE_URL="https://sua-url-do-projeto.supabase.co"
SUPABASE_KEY="sua-chave-anonima-ou-service-role"
API_PORT=":8080"
SECRET_KEY="sua-chave-secreta-base64-aqui"
```

> **Nota:** Nunca compartilhe o arquivo `.env` real em repositórios públicos.

## 3. Configuração do Banco de Dados

Atualmente, a string de conexão com o banco de dados (PostgreSQL/Supabase) está definida no código.

> **⚠️ Atenção:** Para alterar o banco de dados utilizado, edite a variável `ConnectionString` no arquivo:
> `config/config.go`

Caso contrário, a aplicação tentará conectar no banco de dados padrão definido.

## 4. Instalação das Dependências

Abra o terminal na pasta raiz do projeto (`HareID`) e execute o comando abaixo para baixar todas as bibliotecas necessárias:

```bash
go mod tidy
```

Isso garantirá que pacotes como `chi`, `pgx` e `godotenv` estejam instalados.

## 5. Executando a API

Para rodar a aplicação corretamente (garantindo o carregamento correto do `.env`), **você deve executar a partir da pasta `cmd/api`**.

Siga estes passos no terminal:

1.  Navegue até a pasta do executável principal:
    ```bash
    cd cmd/api
    ```

2.  Execute a aplicação:
    ```bash
    go run .
    ```
    *(Ou: `go run main.go router.go application.go`)*

Se tudo der certo, você verá logs como:
```
Db Conn OK
application started at port: :8080
Swagger UI: http://localhost:8080/swagger/index.html
```

## 6. Acessando a Documentação (Swagger)

Com a API rodando, acesse a documentação interativa para testar as rotas:

*   **Link**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
