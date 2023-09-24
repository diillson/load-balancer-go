# **Construindo um Load Balancer com Go**

O balanceamento de carga é uma técnica essencial para distribuir tráfego de entrada de maneira eficiente entre vários servidores back-end. Isso não apenas melhora a latência e a capacidade de processamento do servidor, mas também garante uma distribuição equitativa do tráfego. Neste artigo, exploraremos a construção de um load balancer simples usando a linguagem de programação Go e diversas bibliotecas complementares.

## O que a aplicação entrega

Esta aplicação fornece um load balancer simples que direciona solicitações de entrada para vários servidores back-end. Ele oferece funcionalidades como:

* Distribuição de tráfego baseada no número de conexões ativas.
* Adição e remoção dinâmica de servidores.
* Verificações de saúde para garantir que as solicitações sejam direcionadas apenas para servidores saudáveis.
* Listagem de todos os servidores disponíveis.

## Funcionalidades Principais

### 1. Proxying

O coração do load balancer. Com base na saúde e no número de conexões ativas, direciona a solicitação do cliente para o servidor mais adequado.

### 2. Gerenciamento de Servidores

Permite a adição e remoção dinâmica de servidores. Isso é útil para manter o balanceamento eficiente e adaptável.

### 3. Verificações de Saúde

Os servidores back-end podem ter problemas. A verificação de saúde assegura que apenas servidores saudáveis recebam tráfego.

### 4. Listagem de Servidores

Uma visão geral dos servidores disponíveis, permitindo monitoramento e gestão.

## Bibliotecas e Ferramentas Utilizadas

1. Go (Golang): A linguagem de programação principal usada.
2. Gin: Um framework web HTTP de alto desempenho e fácil de usar para Go. Foi usado para facilitar a criação de rotas e handlers.
3. Viper: Uma biblioteca completa para gerenciamento de configurações em Go, facilitando a carga de configurações de diversos formatos.
4. Logrus: Uma biblioteca estruturada de logging em Go, oferecendo flexibilidade e personalização em logs.
5. Net/HTTP e Net/URL: Bibliotecas padrão do Go para lidar com requisições HTTP e URLs.

## Arquitetura e Design do Código

A aplicação foi desenvolvida considerando boas práticas de programação e design:

* **Arquitetura Modular (ou Standard em como desenvolver código em Go):** Não a foco em uma arquitetura específica XPTO.

* **Separation of Concerns:** O código foi dividido logicamente para facilitar a manutenção e escalabilidade. Por exemplo, a manipulação de requisições HTTP foi colocada no pacote **api**, enquanto a lógica do load balancer reside no pacote **loadbalancer**.

* **Logging Apropriado:** Logrus foi utilizado para fornecer insights valiosos sobre o comportamento da aplicação em tempo real.

* **Configuração Dinâmica:** Viper foi usado para permitir uma configuração flexível, permitindo adicionar ou remover servidores de um arquivo **config.yaml**.

# **Build**

### MacOS
    #amd64
    GOOS=darwin GOARCH=amd64 go build -o cmd/balancer cmd/main.go

    #arm64
    GOOS=darwin GOARCH=arm64 go build -o cmd/balancer cmd/main.go

### Linux

    # amd64
    $ GOOS=linux GOARCH=amd64 go build -o cmd/balancer cmd/main.go

    # arm64
    $ GOOS=linux GOARCH=arm64 go build -o cmd/balancer cmd/main.go

### Windows

    # amd64
    $ GOOS=windows GOARCH=amd64 go build -o cmd/balancer.exe cmd/main.go
    
    # arm64
    $ GOOS=windows GOARCH=arm64 go build -o cmd/balancer.exe cmd/main.go

# **Execução**

    $ cd cmd/
    $ ./balancer

# **Interações**

### Adicionando um novo servidor

    $ curl -X POST -H "Content-Type: application/json" -d '{"URL": "http://example2.com:80", "ActiveConns": 0, "Healthy": true}' http://localhost:3000/addServer

### Listando todos os servidores, informações e status de saúde.

    $ curl -X GET http://localhost:3000/list

### Removendo um servidor

    $ curl -X DELETE -H "Content-Type: application/json" -d '{"URL": "http://example2.com:80"}' http://localhost:3000/removeServer

### Verificando a saúde do LoadBalancer

    $ curl -X GET http://localhost:3000/health

Em resumo, construi um load balancer "simples", mas poderoso, usando Go. A aplicação é eficiente, fácil de configurar e adaptável a mudanças. A combinação de Go com bibliotecas modernas e um design de código cuidadoso torna a solução robusta e pronta para cenários do mundo real. Se você está procurando uma introdução prática ao balanceamento de carga ou deseja entender melhor como implementar soluções de infraestrutura em Go, este projeto serve como um excelente ponto de partida.