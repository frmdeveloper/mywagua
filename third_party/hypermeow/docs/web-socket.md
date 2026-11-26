# Transporte observado no WhatsApp Web

Em 10/11/2026, às 23:43 de São Paulo, um perfil novo do Google Chrome 157.0.8103.51 abriu a página oficial de pareamento. O navegador iniciou dois WebSockets para `web.whatsapp.com/ws/chat`. A conexão na porta 443 registrou quatro frames binários enviados e três recebidos. A conexão na porta 5222 registrou zero nos dois sentidos durante a amostra. O DevTools informou os status 200 e 101, respectivamente. Esses códigos descrevem a negociação observada no navegador. O dispositivo permaneceu despareado. [Observação bruta, sem payloads](web-socket-observation.json).

Para repetir, abra um perfil temporário de Chrome, inicie Network no DevTools, selecione o filtro WS e navegue para a [página oficial](https://web.whatsapp.com/). Registre horário, versão do navegador, host, porta, caminho e contagem de frames por direção. Encerre antes de parear. A porta que transporta os dados depende da rede. A primeira tentativa com Chromium headless recebeu a página de navegador incompatível e não abriu sockets.

## Caminho no SDK

1. `ConnectRace` disputa as aberturas em paralelo. O primeiro sucesso inicia a leitura. Conexões que já abriram e perderam são encerradas, e tentativas pendentes são canceladas. O cliente Go usa a disputa quando a configuração traz duas URLs. [Conexão do cliente](../client.go) e [implementação e testes](../socket/race_test.go).
2. O SDK inicia o handshake Noise e valida a cadeia de certificados antes de passar à leitura dos nós do protocolo. O nível de transporte é distinto da criptografia ponta a ponta entre os dispositivos de remetente e destinatário. A Meta descreve identidades individuais por dispositivo e sincronização de estado entre dispositivos, incluindo favoritos. [Handshake no SDK](../handshake.go) e [arquitetura multidispositivo da Meta](https://engineering.fb.com/2021/07/14/security/whatsapp-multi-device/).
3. O leitor do SDK reconstrói frames internos usando um comprimento de três bytes. Um cabeçalho pode chegar dividido ou vários frames podem chegar no mesmo bloco. O keepalive mede respostas durante a conexão e pode provocar reconexão após falhas repetidas. A resposta HTTP de um comando, o ACK de transporte e o recibo de leitura têm significados distintos. [Leitor de frames](../socket/framesocket.go), [keepalive](../keepalive.go) e [eventos do SDK](../types/events/events.go).

O navegador fornece evidência direta dos endpoints e da troca binária. O comportamento do handshake, o formato interno e a recuperação descritos acima vêm do código Go. A amostra pública não contém conteúdo, confirmações de leitura nem material de sessão.

## Testes reproduzíveis

Execute `go test ./socket -count=1` e `go test ./socket -run TestFrameStreamBoundaries -bench BenchmarkSplitFrameHeader -benchmem -count=3`. O primeiro comando também cobre abertura parcial, cancelamento das tentativas pendentes e encerramento dos concorrentes. O segundo percorre todas as divisões de um fluxo sintético em três blocos e mede um payload de 64 KiB cujo cabeçalho está dividido. Em Windows/amd64 com Go 1.26.10, a remontagem alocava 73.728 bytes por operação. A mudança para um buffer de cabeçalho de três bytes mediu zero bytes alocados nas três repetições. A métrica cobre a remontagem dessa fixture em memória.
