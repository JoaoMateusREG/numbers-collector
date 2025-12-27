
---

```markdown
# Numbers Collector API

API desenvolvida em Go para validação e armazenamento de registros (CPF e número) utilizando banco de dados SQLite.

## 🚀 Como Executar
1. Certifique-se de ter o Go instalado em sua máquina.
2. Execute o comando:
   ```bash
   go run main.go

```

3. A aplicação estará disponível na porta: **7531**

## 📌 Endpoints

### `POST /registro`

Envia um novo registro para o banco de dados.

**Campos obrigatórios no JSON:**
| Campo | Tipo | Descrição |
| :--- | :--- | :--- |
| `cpf` | String | CPF com ou sem formatação (pontos e traços). |
| `numero` | Inteiro (int64) | Número com exatamente 11 dígitos. |

**Exemplo de Payload:**

```json
{
  "cpf": "123.456.789-01",
  "numero": 11988887777
}

```

---

## 🛠 Regras de Negócio

* **Validação de CPF:** A API remove caracteres especiais e valida o CPF através do algoritmo oficial de dígitos verificadores. CPFs matematicamente inválidos são rejeitados.
* **Normalização:** O CPF é salvo no banco de dados apenas como números, garantindo integridade na busca.
* **Atualização Automática (Upsert):** Caso um CPF já cadastrado seja enviado com um novo número, a API **atualizará o registro existente** em vez de criar um novo.
* **Validação de Número:** O campo `numero` deve possuir obrigatoriamente 11 dígitos.
* **CORS:** Configurado para aceitar requisições de qualquer origem, facilitando a integração com front-ends.

## 🗄 Persistência

Os dados são armazenados localmente em um arquivo chamado `dados.db` (SQLite).