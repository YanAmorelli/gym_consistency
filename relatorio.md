# Relatório de Análise de Código: gym_consistency

Este relatório detalha as descobertas de code smells, falhas de segurança e más práticas de codificação identificadas no repositório `gym_consistency`.

## 1. Vulnerabilidades de Segurança Críticas

### 1.1. SQL Injection
*   **Localização:** `handlers/handler_user.go` (`CreateUser`, `LoginUser`, `ForgetPassword`, `ResetPassword`), `handlers/handler_social_media.go` (`RequestFriendship`, `UpdateFriendshipRequest`, `RemoveFriend`).
*   **Descrição:** As queries SQL são construídas usando `fmt.Sprintf` com concatenação direta de dados de entrada do usuário. Isso permite que um atacante insira código SQL malicioso na entrada, executando comandos arbitrários no banco de dados.
*   **Recomendação:** **Prioridade Máxima.** Todas as operações de banco de dados devem usar queries parametrizadas (oferecidas pelo GORM) ou `Prepared Statements` para evitar a concatenação direta de strings.

### 1.2. Armazenamento de Senhas em Texto Plano / Hashing Inseguro
*   **Localização:** `handlers/handler_user.go` (`CreateUser`, `ForgetPassword`, `ResetPassword`).
*   **Descrição:** As senhas dos usuários são armazenadas diretamente no banco de dados (`passwd` coluna) sem hashing ou salting adequados. Adicionalmente, as funções de autenticação (`LoginUser`, `ResetPassword`) dependem de uma função de banco de dados (`auth_login`) que aparentemente lida com senhas em texto plano.
*   **Recomendação:** **Prioridade Máxima.** Implementar um algoritmo de hashing de senhas forte (ex: bcrypt, scrypt, Argon2) com salting para todas as senhas. A autenticação deve ser feita comparando o hash da senha fornecida com o hash armazenado, e não a senha em texto plano.

### 1.3. Fluxo de Redefinição de Senha Inseguro
*   **Localização:** `handlers/handler_user.go` (`ForgetPassword`).
*   **Descrição:** A função `ForgetPassword` gera uma senha aleatória, a atualiza no banco de dados (ainda em texto plano) e a envia por e-mail ao usuário. Isso é extremamente inseguro, pois a nova senha é exposta em trânsito e armazenada de forma vulnerável.
*   **Recomendação:** **Prioridade Máxima.** Implementar um fluxo de redefinição de senha seguro:
    1.  Gerar um token de redefinição único, com tempo de vida limitado.
    2.  Armazenar o hash desse token no banco de dados (não o token em texto plano) associado ao usuário e uma data de expiração.
    3.  Enviar o token por e-mail em um link para uma página onde o usuário define sua nova senha.
    4.  Ao receber a nova senha, verificar o token (hash e expiração) e só então atualizar a senha do usuário (com hashing).

### 1.4. Geração de Senha Aleatória Insegura
*   **Localização:** `services/util.go` (`GenerateRandomPassword`).
*   **Descrição:** A função usa `rand.Seed(time.Now().Unix())` e `rand.Intn()`, que são para geração de números pseudoaleatórios não criptograficamente seguros. Senhas geradas assim são previsíveis e fáceis de quebrar.
*   **Recomendação:** **Prioridade Alta.** Usar o pacote `crypto/rand` do Go para gerar senhas criptograficamente seguras.

### 1.5. `InsecureSkipVerify` em Conexões TLS
*   **Localização:** `services/util.go` (`SendEmail`).
*   **Descrição:** A configuração `d.TLSConfig = &tls.Config{InsecureSkipVerify: true}` desabilita a verificação de certificado TLS para a conexão SMTP, tornando-a vulnerável a ataques "Man-in-the-Middle".
*   **Recomendação:** **Prioridade Máxima.** Remover `InsecureSkipVerify: true`. Configurar corretamente os certificados de confiança no sistema ou no cliente de e-mail. Nunca usar esta opção em ambiente de produção.

### 1.6. CORS Aberto (`AllowOrigins: []string{"*", "*"}`)
*   **Localização:** `main.go`.
*   **Descrição:** A configuração de CORS permite requisições de todas as origens (`*`). Embora útil para desenvolvimento, isso é uma vulnerabilidade de segurança grave em produção, pois qualquer site malicioso pode fazer requisições cross-origin para a API.
*   **Recomendação:** **Prioridade Alta.** Em produção, restringir `AllowOrigins` apenas aos domínios de frontend permitidos pela aplicação.

## 2. Code Smells e Más Práticas de Design

### 2.1. Mistura de Responsabilidades (Handler/Service/Repository)
*   **Localização:** `handlers/handler_social_media.go`, `handlers/handler_stats.go`, `handlers/handler_user.go`.
*   **Descrição:** Os handlers não estão apenas recebendo requisições e enviando respostas; eles contêm diretamente a lógica de negócios, a validação de dados e as operações de banco de dados. Isso viola o Princípio de Responsabilidade Única.
*   **Recomendação:** Refatorar para separar claramente as camadas:
    *   **Handler:** Lida apenas com a requisição HTTP (parsing de input, chamada do serviço apropriado, formatação da resposta HTTP).
    *   **Service:** Contém a lógica de negócios e orquestra operações de persistência.
    *   **Repository (ou Data Access Object):** Encapsula a lógica de interação com o banco de dados (CRUD).

### 2.2. Tratamento de Erros Inconsistente e Exposição de Detalhes Internos
*   **Localização:** Vários arquivos de handler (`handler_social_media.go`, `handler_stats.go`, `handler_user.go`).
*   **Descrição:** As respostas de erro em JSON frequentemente expõem o erro bruto do sistema (`err.Error()`) ao cliente. Isso pode vazar informações sensíveis sobre a infraestrutura ou o banco de dados. Além disso, os códigos de status HTTP não são sempre os mais semanticamente corretos.
*   **Recomendação:**
    *   Criar um formato de resposta de erro padrão (ex: `{ "code": "ERRO_GENERICO", "message": "Ocorreu um erro. Tente novamente mais tarde." }`).
    *   Evitar expor `err.Error()` diretamente ao cliente. Logar os detalhes técnicos internamente.
    *   Utilizar códigos de status HTTP mais apropriados (ex: `400 Bad Request` para validação falha, `401 Unauthorized` para autenticação, `403 Forbidden` para autorização, `404 Not Found` para recursos não encontrados, `409 Conflict` para recursos duplicados).

### 2.3. Validação Básica/Incompleta
*   **Localização:** `services/validation.go` (`ValidateUserData`).
*   **Descrição:** A validação de dados de entrada (`ValidateUserData`) é muito básica, verificando apenas campos vazios e igualdade de senhas. Faltam validações essenciais como formato de e-mail (regex), complexidade da senha (mínimo de caracteres, maiúsculas/minúsculas, números, símbolos) e unicidade do nome de usuário/e-mail.
*   **Recomendação:** Implementar validação abrangente para todas as entradas de usuário. Considerar o uso de uma biblioteca de validação.

### 2.4. Constantes Mágicas / Localização
*   **Localização:** `handlers/handler_social_media.go` (`PENDENT`, `ACCEPTED`, `DECLINED`).
*   **Descrição:** As constantes que representam estados de requisição de amizade estão definidas no arquivo do handler. Embora sejam constantes, sua localização aqui pode não ser ideal se elas precisarem ser usadas em outros pacotes (ex: `models` ou `services`), potencialmente levando a duplicação ou dependências circulares.
*   **Recomendação:** Mover essas constantes para um pacote mais genérico (ex: `models` ou um novo pacote `constants`) onde elas possam ser acessadas por todos os módulos que as utilizam.

### 2.5. Tratamento de `record not found` para `http.StatusInternalServerError`
*   **Localização:** `handlers/handler_stats.go` (`GetDate`).
*   **Descrição:** Para `gorm.ErrRecordNotFound`, a função retorna `http.StatusNoContent`. Em outras partes, um `record not found` poderia levar a um `http.StatusInternalServerError` se o erro não for tratado explicitamente (embora `handler_stats.go` o faça bem para esse caso).
*   **Recomendação:** Padronizar o tratamento para recursos não encontrados. `http.StatusNotFound` (404) é geralmente mais apropriado do que `http.StatusNoContent` (204) quando um recurso específico não pode ser encontrado.

### 2.6. Variável Global de Erro
*   **Localização:** `handlers/handler_stats.go` (`duplicateEntryError`).
*   **Descrição:** `var duplicateEntryError = &pgconn.PgError{Code: "23505"}` define uma variável global para um erro, o que é uma má prática e pode levar a condições de corrida ou comportamento inesperado.
*   **Recomendação:** Evitar variáveis globais para erros. A verificação do código de erro do PostgreSQL deve ser feita diretamente no local de uso: `var pgErr *pgconn.PgError; errors.As(err, &pgErr) && pgErr.Code == "23505"`.

## 3. Melhorias e Refatorações Gerais

### 3.1. Reutilização da Validação de Token (Middleware)
*   **Localização:** Repetição da lógica de `services.ValidateUserToken` em vários handlers.
*   **Descrição:** A lógica para validar o token do usuário é repetida em muitas funções de handler.
*   **Recomendação:** Criar um middleware de autenticação no Echo que valida o token e injeta as `claims` do usuário no contexto da requisição. Isso simplifica os handlers e centraliza a lógica de autenticação.

### 3.2. Padronização do Logging
*   **Localização:** Uso inconsistente de `log.Println` e `log.Error` em diferentes arquivos.
*   **Descrição:** O estilo e o nível de detalhe do logging de erros variam entre os arquivos.
*   **Recomendação:** Adotar uma estratégia de logging consistente para toda a aplicação, definindo níveis de log (INFO, WARN, ERROR, DEBUG) e formatos.

### 3.3. Otimização da Geração de Senha
*   **Localização:** `services/util.go` (`GenerateRandomPassword`).
*   **Descrição:** A lógica de adição de caracteres especiais na senha gerada é um pouco ineficiente e complexa.
*   **Recomendação:** Refatorar a função para uma abordagem mais limpa e eficiente que garanta a presença de diferentes tipos de caracteres (maiúsculas, minúsculas, números, especiais) de forma aleatória e segura.

### 3.4. Refatoração de `UpdateFriendshipRequest`
*   **Localização:** `handlers/handler_social_media.go` (`UpdateFriendshipRequest`).
*   **Descrição:** A lógica de atualização do `request_status` parece ter um erro, utilizando o status atual para a atualização. Operações que envolvem múltiplas escritas no banco de dados não estão em transações.
*   **Recomendação:** Corrigir a lógica para garantir que o status correto seja aplicado. Encapsular todas as operações de escrita em uma transação de banco de dados para garantir atomicidade.

### 3.5. Tratamento de `ioutil.ReadAll`
*   **Localização:** `services/google_validation.go` (`getGooglePublicKey`).
*   **Descrição:** Uso de `ioutil.ReadAll`, que foi depreciado no Go 1.16.
*   **Recomendação:** Substituir por `io.ReadAll`.

### 3.6. Tratamento Incorreto de Erro em `json.Unmarshal`
*   **Localização:** `services/google_validation.go` (`getGooglePublicKey`).
*   **Descrição:** Se `json.Unmarshal` falhar, a função retorna `"", nil`, o que pode levar a falhas silenciosas.
*   **Recomendação:** Retornar `"", err` para propagar o erro de unmarshalling corretamente.

---

Este relatório serve como um guia para melhorias no código. Abordar as vulnerabilidades de segurança deve ser a prioridade máxima, seguida pela refatoração para um design mais limpo e manutenível.

**Nota:** Esta seção foi adicionada para criar um commit distinto e permitir a abertura do Pull Request.