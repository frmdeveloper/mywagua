# Plano executavel: Cloud Go-only e HyperMeow WPPConnect

Inicio: 2026-11-10. Responsavel: joaosouz4dev. Atualizar os checkboxes somente com evidencia.
[x] significa concluido e verificado no escopo descrito; publicacao e testes reais tem itens separados.
Nao declarar 100% enquanto houver requisito ou validacao pendente. Preservar trabalho local existente.
Os registros abaixo sao cronologicos: pendencias de um marco antigo podem estar encerradas nos marcos seguintes. O checklist e o ultimo registro indicam o estado atual.

## Fase 0 - Base e rastreabilidade
- [x] Comparar Go, Node e imagem Zuck v2.4.2 por revisoes/digests fixos.
- [x] Registrar autorizacao do usuario para publicar e fonte publica indicada: polymorfa/hypermeow.
- [x] Criar checkouts isolados Cloud/Manager/Go; preservar checkouts humanos.
- [x] Criar checkout isolado worker a partir da main 3ca6010.
- [x] Inventariar fonte publica HyperMeow, licencas, relacao com tulir e diferencas da imagem.
- [x] Confirmar se a imagem contem fonte privado completo; nao equiparar binario com fonte recuperado.

## Fase 1 - Cloud e pools somente Go
- [x] Cloud API e Manager rejeitam selecao diferente de Go; defaults e retorno coerentes.
- [x] Reservas do Manager filtram Node mesmo com configuracoes antigas.
- [x] Worker gera apenas Go, sem reposicao de Evolution/Wuzapi/Node por env legado.
- [x] Worker rejeita novas criacoes de outros providers antes de Docker/DB.
- [x] Preservar operacao/limpeza de sessoes existentes sem reatribuir volumes ou credenciais.
- [x] Testar defaults, env legado, pedido forjado, concorrencia e falhas de provisionamento.
- [x] Identificar ambiente em uso, configurar reserva Go, publicar fontes testadas e verificar pool real.

## Fase 2 - Documentacao completa sem escolha de motor
- [x] Exemplos de criacao com JSON minimo, sem provider/runtime obrigatorios.
- [x] Guia publico: autenticacao/scopes, QR, status, texto, midia, lifecycle, limites e erros.
- [x] Corrigir rotas/payloads/eventos documentados conforme API real; distinguir recurso disponivel e pendente.
- [x] Publicar referencia operacional: pool, imagem, persistencia, backup, logs, rollout/rollback.
- [x] Atualizar traducoes PT/EN e verificar links/exemplos/build.

## Fase 3 - Projeto publico HyperMeow
- [x] Preparar repositorio wppconnect-team/hypermeow a partir de polymorfa/hypermeow, com historico/licencas.
- [x] Documentar origem, escopo SDK vs servidor Zuck e modulos/imports para uso no Go.
- [x] Rodar testes/vet/build e CI da base; registrar qualquer dependencia externa indisponivel.
- [x] Publicar repositorio publico e conferir visibilidade/SHA/licencas.
- [x] Integrar atualizacoes de tulir preservando extensoes HyperMeow, com conflitos/testes resolvidos.
- [x] Validar adaptacao do server-go ao SDK antes de trocar dependencia em producao.

## Fase 4 - Sincronizacao diaria
- [x] Agendar comparacao diaria com tulir/whatsmeow as 09:30 America/Sao_Paulo.
- [x] Executar em checkout isolado, resolver principal no provedor e registrar SHAs de origem/destino.
- [x] Preservar patches e licencas; nunca usar reset/push force para igualar forks.
- [x] Testar, incorporar e enviar atualizacoes compativeis; falhas ficam registradas sem promover.
- [x] Executar primeira comparacao agora e validar configuracao da tarefa agendada.

## Fase 5 - Integracao Cloud e confiabilidade
- [x] Corrigir webhook por sessao na ativacao do pool, inclusive webhookUrl/webhook, restart e isolamento.
- [x] Implementar HMAC opcional nos webhooks globais do Go, com testes e documentacao.
- [x] Idempotencia persistida de envio, recuperacao de resposta e bloqueio de reenvio ambiguo; incerteza do protocolo continua explicita.
- [x] Falhas de entrega: consulta por sessao/destino/estado, replay seletivo auditado e metricas de backlog/capacidade, com testes de falha e recuperacao.
- [ ] ADIADO pelo usuario em 10/11/2026: nenhum monitoramento ou destino de notificacao por ora. Nao instalar nem disparar alertas operacionais; as regras seguem como exemplos e nao ha entrega validada.
- [x] Conta de teste real: QR, reconexao, texto/arquivo, ACK de entrega/leitura, resposta recebida, idempotencia e download de midia depois de restart. Outro workspace recebe 404.
- [ ] Logout e pareamento por codigo de telefone separados do QR.
- [ ] ADIADO pelo usuario em 10/11/2026: jornadas com segunda conta pareada e validacao Android/iPhone/Web. Nao bloqueiam a implementacao automatizada; nao contam como aprovadas.

## Fase 6 - Paridade funcional Zuck
- [x] Converter inventario de 303 operacoes em matriz executavel de contratos e criterios por jornada.
- [x] SDK: incorporar recursos publicos HyperMeow e identificar extensoes privadas nao presentes.
- [ ] Adaptadores de contatos, grupos/comunidades, newsletters, status e midia.
- [x] J10: corrigir espera/cancelamento de ACK de visualizacoes no SDK antes de expor a operacao.
- [x] J10 no server-go (contratos locais com cliente de teste): convites, reacoes, visualizacoes e atualizacoes de mensagens com contratos estritos, isolamento e idempotencia nas mutacoes.
- [x] J10: protocolo de ACK no SDK, fixtures HTTP/SQLite no Go, OpenAPI e guia nativo; CI aprovado e fontes incorporadas.
- [x] J10 na Cloud: escopos separados para leitura, escrita e administracao, resultados isolados, documentacao PT/EN e testes de navegador/build/CI.
- [x] J10: publicar a imagem no worker e renovar/verificar as tres reservas em Go; smoke pelo Manager sem telefone pareado.
- [x] J10: Cloud 5b9f3e8 publicada; /docs em portugues com o contrato de canais, quatro rotas protegidas responderam 401 e rota desconhecida 404.
- [ ] J10: aceite de owner/subscriber em duas contas e nos clientes Android/iPhone/Web adiado pelo usuario.
- [x] Redis Streams: rejeicao/recuperacao, trim/exact maxlen e identidade do evento com Redis real no CI; incorporado em [Go #22](https://github.com/wppconnect-team/wppconnect-server-go/pull/22).
- [x] SQS: destino opcional, rejeicao de fila inexistente, recuperacao e deduplicacao FIFO com ElasticMQ; incorporado em [Go #24](https://github.com/wppconnect-team/wppconnect-server-go/pull/24). Sem teste de credenciais/IAM/KMS AWS reais.
- [x] n8n: importar workflow e credencial temporarios, publicar, recusar segredo errado e processar evento autenticado em n8n 2.47.3 isolado ([CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42902244064/job/130130987174)).
- [ ] Integracao Chatwoot e armazenamento automatico de midia. CRUD/segredos de brokers por sessao (J15) continuam separados dos adaptadores globais opcionais.
- [x] J11 no server-go (contrato local): catalogo proprio, produto e colecoes, limite e cursor, escopo da conta e preco sem conversao.
- [x] J11: fontes de leitura incorporadas em [Go #23](https://github.com/wppconnect-team/wppconnect-server-go/pull/23), CI com race e imagens.
- [x] J11: leitura publicada na Cloud, protecao das tres rotas, guia EN/PT e tres reservas Go saudaveis; [recibo](docs/brokers-catalog-acceptance.json).
- [ ] J11: aceite de catalogo em conta Business com produtos/colecoes reais.
- [ ] Comercio: catalogos/produtos/pedidos/pagamentos com contratos e idempotencia.
- [x] J05/Z126: contrato HTTP experimental de carrossel PNG/JPEG com respostas rapidas; validacao antes de upload, isolamento, deduplicacao e payload protobuf testados; CI 42901305610 aprovado e main 93304d79bd09af36386035170e78886e6b151e56.
- [ ] Interativos modernos: botoes/listas/carrosseis/flows com teste Android/iOS/Web.
- [ ] Chamadas: sinalizacao/audio/video/grupos, autenticacao de streams e limites de recursos.
- [ ] Testes de protocolo para cada recurso; ausencia de fonte ou teste nao conta como paridade.

## Fase 7 - Aceite e encerramento
- [x] CI, race, integracoes PostgreSQL/brokers/S3, build e smoke multiarch aprovados para as mudancas incorporadas ate Go e01698b e Cloud 35ac0c4.
- [x] Restauracao offline do backup consistente da conta de teste: integridade dos cinco bancos SQLite, sem abrir uma segunda conexao WhatsApp.
- [x] Ensaio isolado de recuperacao idempotente com 1 e 4 clientes: 72 leituras autenticadas identicas, quatro registros preservados depois de restart e quatro consultas entre sessoes negadas.
- [ ] Carga de protocolo, reconexao, midia e backlog com metas de capacidade aceitas; o microbenchmark HTTP nao cobre esses perfis.
- [x] Mudancas implementadas publicadas a partir das fontes verificadas; health/logs/versao confirmados em 10/11/2026.
- [ ] Checklist integral concluido com links de evidencia. Somente entao declarar 100%.

## Evidencias
- Analise: server-go/docs/research/server-go-parity-2026-11-10.md.
- Primeira implementacao: IMPLEMENTATION.md.
- Base publica/privada: public-hypermeow-foundation.md (concluido).
- Este plano e o registro principal; complementar com recibos de testes/publicacao sem incluir segredos.

## Registro de validacao - 2026-11-10
- HyperMeow publico: https://github.com/wppconnect-team/hypermeow, main bf5bab849ea6f666f490e32487cec0a7c84dbeb5. Historico e licencas preservados.
- CI HyperMeow aprovado: https://github.com/wppconnect-team/hypermeow/actions/runs/42860662627 (Go antigo/atual, testes, build, WASM, MEX e pre-commit).
- Automacao ativa: sincronizar-hypermeow-wppconnect-diariamente, diariamente 09:30 America/Sao_Paulo. Configuracao conferida pela ferramenta do aplicativo. As regras da fase 4 estao no prompt; a primeira integracao ainda tem 15 conflitos a resolver.
- Auditoria das 10 camadas AMD64: imagem nao contem fonte privado completo; contem binario, UI e Swagger. O projeto publico nao e uma recuperacao desse fonte.
- Cloud: 43 testes locais e build Next.js aprovados. PR https://github.com/joaosouz4dev/wppconnect-cloud/pull/37.
- Worker: 13 testes, 3 smokes e build TypeScript aprovados. PR https://github.com/joaosouz4dev/wppconnect-cloud-manager-worker/pull/6.
- Manager: 14 testes unitarios aprovados; integracoes locais nao executaram por falta de PostgreSQL (ECONNREFUSED). CI com banco isolado pendente. PR https://github.com/joaosouz4dev/wppconnect-cloud-manager/pull/5.
- Guia operacional esta versionado nos PRs; publicacao no ambiente Cloud continua pendente na fase 1.

## Publicacao Cloud Go-only - 2026-11-10
- Cloud main 70c00e087ccc12d0235fb0ee0ab58dbb7f16dec7, deployment dpl_QKxcN2iGNJZ1LF2Xq6raKu9JWfVg. Dominio confirmado; /docs HTTP 200, criacao anonima HTTP 401. Evidencia cloud-release.json.
- Manager main 3ade44f97868c008c30a1e0894f9636cfa851c51; worker main f87a9a9b9d18579f5f1aa8521a125c75a8c3e9cb. CI aprovado, inclusive PostgreSQL, concorrencia e rollback. Imagens e versoes anteriores em go-only-services-release.json.
- Canary: rejeicao de criacao nao-Go no worker/Manager, autenticacao, bancos e proxy da sessao existente aprovados. Containers temporarios removidos; evidencia go-only-service-canary.json.
- Producao: WPP_RUNTIME_POOLS=go:go=3, tres reservas Go, pedido Node HTTP 422, sessao antiga CONNECTED. Outros containers do host preservados.
- Smoke real com workspace isolado: criacao minima HTTP 200/Go, QR PNG valido, outro tenant HTTP 404, cota HTTP 409, exclusao de teste concluida e reposicao automatica para tres reservas Go. Nao houve envio de mensagens nem pareamento de conta de cliente.
- Integracao HyperMeow com tulir em PR https://github.com/wppconnect-team/hypermeow/pull/1; 15 conflitos resolvidos e testes locais aprovados; CI/race pendentes. Ainda nao promovida a main.

## SDK e contratos - 2026-11-10
- HyperMeow sincronizado com tulir c88388f: PR #1 incorporado, merge aae5d9a62408c202fb4dd2f91394153c054a110b. CI Go 1.26/1.27, race, WASM, MEX e pre-commit aprovados nas runs 42866085986 e 42866095663.
- HMAC server-go: PR #15 incorporado, merge 36e41bc; CI 42865161731 aprovado com race, brokers/S3 e smoke amd64/arm64. Ainda nao altera a imagem 0.1.11 da pool.
- Adapter HyperMeow server-go: PR #16 incorporado, merge dc8d124. CI 42866854308 aprovado com race, integracoes e imagens das duas arquiteturas. Modulo remoto fixado em 0987e7d211c8.
- Teste SQLite reproduziu duplicate column na migracao da base anterior v16. Corrigido no SDK e validado com preservacao de chaves Noise/identidade/ADV, segredo WASA, indice faltante e upgrade repetido.
- Matriz estruturada com 303 operacoes e 20 jornadas: parity-contract-matrix.md / .json. Versao publica compacta lossless 482233 bytes, validada por expansao e comparacao integral. Todas as operacoes seguem functionalParityValidated:false no baseline pesquisado; nao sao testes de mensagens reais.
- Webhook por sessao: implementacao em branches server-cloud-webhooks e worker-webhooks, com receiver dedicado, segredo e restricao de egress HTTPS publico; rollout e verificacao ponta a ponta pendentes.

## Webhooks publicados e verificados - 2026-11-10
- Server-go PR #17 e worker PR #7 incorporados apos CI aprovado. Go main 0878a10e51beae17fa199e2abd9b76ddf0f59ac1; worker main 95aa5d279778e84707b5bfde886cdd455bee48fe.
- Imagens construidas de arquivos git de checkouts limpos e hashes conferidos. Go sha256:e249ddc5d0f6a06d65cce3e63e178cc25b68cbf597b0bf41871b869430a21e8a; worker sha256:031ccfa74a61b0e176210d0cbefafdf4e5f7b21dca363b4ee392eea953605fef.
- Canary real: dois receptores HTTPS com chaves distintas, assinatura validada sobre corpo original, primeira entrega 503 seguida de retry com ID estavel, zero assinaturas invalidas, nenhum evento cruzado, restart preservado. Producao e pool preservadas. Evidencia webhooks-canary.json.
- Rollout: GO_SESSION_WEBHOOKS_ENABLED=true e imagem nova para futuras criacoes. Tres reservas anteriores e sessao Node CONNECTED preservadas. Backup privado registrado em webhooks-release.json.
- Manager real: criacao com webhook em 5,91 segundos, QR PNG, webhook assinado, outro workspace 404, cota 409, exclusao concluida. Nenhum telefone pareado. Evidencia webhooks-live-smoke.json.
- Receptor HTTPS temporario, rota e segredos de teste removidos. Evidencia webhooks-cleanup.json.
- Cloud PR #38 incorporado. Documentacao PT/EN publicada em https://cloud.wppconnect.io/docs; HTTP 200 e exemplo webhookSecret confirmado. Deployment dpl_Bdf6d87gR9EaGcUjEtmnizxoQ73j, fonte e895a7d. 43 testes, build e CI aprovados.
- HyperMeow PR #2 incorporado: matriz publica completa e testes herdados restaurados, CI Go antigo/atual aprovado.
- Proxima etapa em andamento: ledger persistido de comandos de envio, com chave por sessao, recuperacao de resposta e bloqueio de reenvio ambiguo. Branch server-idempotency; testes focados aprovados, suite completa/CI/rollout ainda pendentes.

## Criterios de fechamento das lacunas restantes
- Idempotencia: concorrencia, reinicio, falha de persistencia, timeout parcial, isolamento, rejeicao em imagens antigas e transporte do header por todos os proxies; depois smoke no artefato publicado.
- DLQ: listar entregas falhas por sessao/destino, motivo sanitizado, replay seletivo sem repetir destinos concluidos, auditoria, metrica de idade/ocupacao e ensaio de indisponibilidade/recuperacao.
- Paridade: implementar por jornada da matriz, com fixtures de protocolo e contrato; comparar resultado observavel e nao somente nome de rota. Marcar cada operacao somente depois de validar comportamento.
- WhatsApp real: destinatario autorizado recebido; depende agora de conta de teste pareada. Nunca usar sessoes de clientes para envios de teste.
- Atualizacao de aceite: destinatario autorizado recebido em 10/11/2026 e mantido em registro privado, sem publicar o telefone. Instancia Go isolada criada na imagem validada; QR apresentado ao usuario. Pareamento ainda pendente, nenhuma mensagem real enviada.
- Fonte privada: todas as camadas da imagem foram inspecionadas e nao contêm o fonte completo. As extensoes publicas foram preservadas; recursos exclusivamente binarios precisam de implementacao propria ou fornecimento do fonte. Nao declarar clonagem integral.

## Confiabilidade incorporada - 2026-11-10
- Server-go PR #18: ledger persistido de envio com chave por sessao, fingerprint do corpo/rota/Content-Type, resposta recuperavel, reserva antes do efeito de protocolo, bloqueio de reenvio em resultado incerto e limites de armazenamento. CI 42873726008 aprovado, incluindo race e imagens multiarch. Merge c29f4ac.
- Server-go PR #19: filtros de entregas por destino/estado, motivo sanitizado, retry somente de destino falho, auditoria transacional, migracao da base anterior e metricas de idade/capacidade. CI 42875435897 aprovado. Main 514347cd5e7dd69348a1ef1d63e51f31ae29cf14.
- Cloud PR #39: preservacao de Idempotency-Key, consulta send-result e verificacao previa do ledger para recusar imagens antigas sem enviar. 48 testes, build e CI aprovados. Main 35ac0c4080dd375a3efa1625567e55d93c856958; deployment dpl_ATXHSjh4bJDtWKBKoQC5CrvT4M6R publicado e /docs conferido HTTP 200.
- Imagem Go sha256:e2c292c3259007d8d74c72422b847510ab2c75d43692e5dfb265345523dcb2ec, construida da main limpa e arquivo-fonte SHA256 c0c105975fab5eb167cedfdf2fac1b124480d2c775e7179846d283e65a8b86c3.
- Canary real da imagem: falhas HTTP repetidas, estado failed/retry_exhausted, recuperacao seletiva com assinatura valida, auditoria, bloqueio entre sessoes, metricas, resultado idempotente, conflito de payload e persistencia apos restart. Nenhuma mensagem enviada e nenhum telefone pareado. Fixtures e rede temporaria removidas; producao preservada.
- Primeiro ensaio abortou durante conexao inicial reiniciada no boot. O script foi corrigido para aguardar esse estado transitorio; ensaio completo seguinte aprovado. Nao foi uma falha funcional ocultada.
- Regras Prometheus foram documentadas como exemplos, sem declarar monitoramento instalado. Essa ativacao permanece pendente e separada da implementacao de metricas.
- Rollout e renovacao das tres reservas concluidos via alocacoes de teste proprias, sem apagar reservas diretamente. Todas usam a imagem e2c292c3. Smoke pelo Manager confirmou transporte da chave, recuperacao de resposta e conflito de payload; alocacoes de teste excluidas. Recibo reliability-live-pool.json.
- Auditoria final 18:46 UTC: Manager/worker em execucao, tres Go healthy, nenhum OOM, panic/fatal/falha de persistencia nos dez minutos consultados, sessao antiga CONNECTED e containers temporarios ausentes. Recibo reliability-final-audit.json.
- Esses smokes usaram contas nao pareadas e payloads invalidos para nao enviar mensagens reais. O protocolo de envio foi exercitado com clientes simulados nos testes; recebimento real em telefone permanece pendente.

## Sequencia da fase de paridade funcional
1. Contatos, grupos/comunidades e newsletters: partir dos adaptadores ja existentes e da matriz, fechar diferencas de payload, ownership, paginacao e eventos. Aceite com fixtures de protocolo e duas contas de teste.
2. Midia e integracoes: persistencia automatica e retencao, Chatwoot/n8n, Redis Streams/SQS; testar indisponibilidade, recuperacao, duplicatas e limites antes de publicar cada adaptador.
3. Comercio: leitura de catalogo/produto antes de mutacoes e pedidos. Preservar moeda/unidade, diferenciar aceite de protocolo de pagamento confirmado e usar idempotencia nas mutacoes suportadas.
4. Interativos: validar botoes/listas, depois carrosseis/flows, com fixtures e recebimento observado em Android, iOS e Web. Nao converter resposta HTTP aceita em paridade comprovada.
5. Chamadas: separar rejeicao/sinalizacao de transporte de audio/video e chamadas de grupo. Exigir autenticacao de streams, limites, limpeza, codec compativel e teste real de ponta a ponta.
6. Carga e restauracao: medir CPU/RSS/disco, p95/p99, backlog e reconexao com perfis reproduziveis; restaurar backup consistente e repetir jornadas antes do aceite final.

Cada jornada deve atualizar a matriz com SHA, contrato, teste, resultado e lacuna restante. As 303 operacoes do baseline continuam identificadas como baseline, sem substituir evidencia funcional por uma contagem de rotas novas.

## Pareamento real - QR expirado
- [x] Usuario informou "QR invalido ou expirado"; diagnostico confirmou instancia saudavel e ciclo anterior encerrado sem CONNECTED.
- [x] Substituida imagem estatica do chat por pagina HTTPS temporaria, caminho privado aleatorio, validade de 30 minutos, sem cache e sem expor credencial do Manager.
- [x] Verificados PNG atual e mudanca real da versao do QR. A pagina consulta o estado a cada segundo e reinicia o ciclo encerrado.
- [x] Validado JavaScript efetivamente renderizado; erro de escape identificado e corrigido antes de entregar o link. Teste test-live-pairing.cjs aprovado para rotacao, ocultacao de QR antigo, estado conectado e expiracao.
- [x] Confirmar CONNECTED depois do escaneamento: servidor confirmou conexao em 10/11/2026 19:19 UTC; usuario confirmou sucesso.
- [x] Executar envio/recebimento reais autorizados. Primeiro envio retornou 502 e foi preservado no ledger sem reenvio automatico; apos a correcao, novo teste identificado de texto e arquivo recebeu HTTP 201, ACK de entrega/leitura e resposta real.
- [x] Removidos helper HTTPS, rota e configuracao privada depois de confirmar CONNECTED. Timer de limpeza desativado; sessao de teste preservada. Recibo live-pairing-cleanup.json.
- [x] Agendada limpeza automatica do helper expirado via timer systemd dedicado; limpeza antecipada apos pareamento concluida e timer desativado. Script verifica nome, label e diretorio antes de remover somente o helper e sua configuracao.

## Correcao encontrada no aceite real
- [x] Isolar falha de resolucao: probe sem envio confirmou GetUserInfo com perfil valido sem LID, ResolveLID recusando destinatario e IsOnWhatsApp retornando contato registrado em LID. Conta de teste diferente do destinatario; nenhum dado privado publicado.
- [x] SDK: fallback para consulta de contatos quando perfil omite LID; exige correspondencia da consulta e um unico contato registrado em LID. Preserva device e mappings canonicos.
- [x] Regressao reproduzida antes da correcao por serializacao e parsing reais de USync. Seis cenarios, suite completa e vet aprovados; CI Go antigo/atual, race, WASM, MEX e pre-commit aprovados na run 42886170997.
- [x] HyperMeow PR #4 incorporado: d326350a86e98c405c9f6adf14cb6bd03ea8e66a. Fonte da correcao 44b882797ce70f44ce41d1d31d7461824d92a6c0.
- [x] Atualizar dependencia do server-go, testar, incorporar e publicar imagem validada; repetir aceite real e renovar reservas. Server-go PR #20, main e01698bf3e8b0b9d9dc749e7ffbb49e65deab22f; CI 42887160468 aprovado.

## Aceite real e pool atualizada - 2026-11-10
- Imagem: sha256:b0d96f63e775038ff4c71cfdf4630763712fe8e978ff555dedb685cae8ed9b5b. Arquivo-fonte da main limpa: SHA256 17c5ed3a0018f94f2b3527b547ccc751b90092060067d62e402d43cb454e735e.
- Canary isolado repetiu recuperacao de webhook, auditoria, isolamento, ledger e restart. Fonte conferida novamente antes do rollout. Evidencias recipient-fix-built.json e recipient-fix-canary.json.
- Primeiro ensaio de troca da conta de teste falhou por alias configurado na rede bridge padrao. Script corrigido; imagem anterior restaurada com o volume preservado e CONNECTED confirmado. Segunda troca concluida com backup consistente, mesma porta/configuracao e pareamento preservado. Nao houve alteracao de sessoes de clientes.
- Texto e arquivo entregues e lidos pelo destinatario autorizado; resposta exata de teste recebida pelo servidor. Repeticao da mesma chave recuperou resultado, alteracao de payload retornou 409, outro workspace nao acessou o resultado. Recibo publico sem identificadores: phone-acceptance-summary.json.
- Restart preservou pareamento, os dois resultados bem-sucedidos e a falha inicial. Download do arquivo depois do restart identico byte a byte. Backup restaurado offline com integridade dos cinco bancos SQLite validada. Recibo phone-restart-validation.json.
- Tres reservas anteriores renovadas por alocacoes de teste proprias e exclusao pelo Manager. Todas usam a nova imagem; Manager e worker saudaveis, sessao anterior preservada. Recibos recipient-fix-release.json, recipient-fix-live-pool.json e recipient-fix-final-audit.json.
- Helper de pareamento, rota HTTPS, configuracao privada e container do probe removidos. Sessao Go de teste permanece pareada para as jornadas restantes. Fonte da imagem privada, logout, segunda conta, alertas, carga e paridade integral continuam com as limitacoes registradas acima.


## Canais e integracoes - 2026-11-10

- [HyperMeow #6](https://github.com/wppconnect-team/hypermeow/pull/6): cancelamento, desconexao e ACK negativo nao produzem sucesso de visualizacao. CI 42893049289; main 099af656ac43dcd4b3e2688f8a880ee76e2c4e75.
- [Server Go #21](https://github.com/wppconnect-team/wppconnect-server-go/pull/21): oito operacoes, seis mutacoes com ledger persistente e duas consultas separadas de resultado. Testes de contrato HTTP, isolamento, permissao e replay; CI 42895236212 com race, integracoes e imagens amd64/arm64. Main 4fb3d7f13de475af0e6fa7bcd89b733b12105121. [Guia de canais](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/newsletters.pt-BR.md).
- [Cloud #40](https://github.com/joaosouz4dev/wppconnect-cloud/pull/40): 49 testes, controle de escopos, typecheck, build e CI 42895717043. Chromium: guia em ingles (1440 px) e portugues (375 px), link do contrato, foco de teclado e largura da secao conferidos.
- [Server Go #22](https://github.com/wppconnect-team/wppconnect-server-go/pull/22): adaptador Redis Streams opcional, fixture de Redis real no CI, inspecao e replay por destino. Nenhuma fila ou stream de cliente configurado.
- A matriz original continua sendo o retrato de 303 operacoes da versao comparada. Nomes e payloads sao nativos; Z025 nao tem filtro before no SDK. A reacao recebe 202 sem ACK remoto e um convite administrativo nao concede o papel.

- Canary do binario e rollout das tres reservas Go aprovados: SHA da imagem 59f610bd5d2060cd4b1eecff4ac134f7a41a49d696e97bdac68461b4eeff508a. Auditoria 20:41 UTC: tres reservas healthy, containers temporarios removidos e nenhum OOM, panic, fatal ou falha de persistencia nos logs verificados. A conta de testes pareada e a alocacao preexistente foram preservadas.
- [Server Go #23](https://github.com/wppconnect-team/wppconnect-server-go/pull/23): leitura do catalogo proprio, produtos e resumos de colecoes. Testes locais de escopo, preco/moeda, visibilidade falsa, cursor, limites e erros aprovados.


## Brokers, comercio e carga - 2026-11-10

- Redis Streams: main ba6dee6414e7630e18ced46c0be12b74fbf69b62, CI 42896819299. O ID de evento permanece estavel mesmo quando uma entrega repetida gera outro ID de entrada no stream.
- Catalogo proprio: main 1c42366e18125015a9d30d7abf3a5cafedeccb31, CI 42898388266. A Cloud usa commerce:read em tres consultas. Testes de conta Business real, mutacoes, pedidos e pagamentos continuam abertos.
- SQS: main 666afa223c12efefcb9ac90a2d6920d9ea3a1b11, CI 42898825516. O teste de integracao usa ElasticMQ, nao uma conta AWS. Tanto SQS como Redis permanecem desativados por padrao.
- [Carrossel experimental #25](https://github.com/wppconnect-team/wppconnect-server-go/pull/25): validacao local HTTP/SQLite/SDK, rejeicao antes de upload e resposta idempotente. Nao enviar para validar clientes fisicos enquanto esse teste estiver adiado.
- Benchmark isolado: limite 0,75 CPU e 256 MiB, 24 consultas sequenciais e 48 com quatro clientes. P95 de 98,50 ms e 575,48 ms respectivamente. Primeiro ensaio nao concluiu o restart por reutilizar a porta efemera anterior; o ensaio completo consultou novamente a porta e recuperou os quatro registros. Sem mensagens de WhatsApp. [Medicoes](docs/benchmark-result-lookup.json) e [script](docs/benchmark-result-lookup.py).

- Publicacao de brokers/catalogo: imagem 49a88afea20f64521ed3bb453aeb4ef3bbc77ba45b43d011abd0645cabf8db78. Tres reservas Go healthy, sem OOM ou erros criticos recentes. Cloud #41, main 3751275b5385c92ea25f5db623669ac97aad4baa, deploy confirmado e guia publico conferido em dois idiomas.

- CI da fonte combinada brokers/catalogo: [42900118182](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42900118182). O OpenAPI no Git era identico byte a byte ao gerado; o primeiro checkout Windows materializou CRLF. Repeticao com os bytes do Git passou e a regra de LF foi incorporada no PR #25.

## Continuidade do aceite

Os criterios de fechamento por entrega estao em [parity-remaining.md](docs/parity-remaining.md). O [validador do progresso](docs/check-parity-progress.py) conferiu os 12 contratos mapeados no SHA 93304d79; uma alteracao proposital removendo o status 202 da reacao foi detectada em copia isolada. Nenhuma equivalencia integral foi inferida.

- [x] [Servidor Go #26](https://github.com/wppconnect-team/wppconnect-server-go/pull/26): destino independente do n8n com autenticacao, falha e replay por evento. No n8n, SQS e Redis o evento de QR informa disponibilidade sem transmitir o codigo de pareamento. [CI com workflow real](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42903287552).
- [x] Imagem 765625f: tres reservas Go saudaveis em 10/11/2026 21:31 UTC, carrossel desligado por padrao, falha e recuperacao assinada verificadas, conexao preexistente preservada. [Recibo sem dados de sessoes](docs/workflows-carousel-acceptance.json).

## Entregas verificadas: consultas, recuperação de fila e guia público

- [x] [Servidor Go #27](https://github.com/wppconnect-team/wppconnect-server-go/pull/27): cache por processo de até 256 verificações bem-sucedidas, por até 30 segundos, sem renovar o prazo a cada acerto. Nova sessão, segredo ou token exige validação. A primeira implementação foi sinalizada pelo CodeQL; a versão incorporada usa HMAC com chave aleatória por processo. [CI da correção](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42905486316).
- [x] [Benchmark comparativo](docs/session-verification-benchmark.json) nas imagens reais de 765625f e fc69f92, mesmo host, mesma fixture, 0,75 CPU e 256 MiB. P95 de 104,34 para 2,63 ms em 24 consultas sequenciais; de 499,75 para 4,05 ms em 48 consultas com quatro clientes. Todas as respostas conferidas byte a byte, quatro registros recuperados após reinício e quatro consultas negadas entre sessões em cada imagem. O resultado é específico da consulta autenticada aquecida; não mede envio no WhatsApp.
- [x] [Pool e canary](docs/session-cache-acceptance.json): imagem b5e6e8113ba9ed0b4ea4bac01b29ad322847b659483bbfb132b60c5b90712c28. Três reservas Go saudáveis, conexão existente preservada, sem OOM ou erros críticos encontrados na janela inspecionada. Auditoria de 10/11/2026 às 21:40 UTC. Os destinos externos opcionais continuam desabilitados na pool.
- [x] [Servidor Go #28](https://github.com/wppconnect-team/wppconnect-server-go/pull/28): teste em processo separado com 128 eventos, duas sessões e dois destinos, SQLite WAL, recusa ao atingir o limite e queda sem fechar o banco. Verifica preservação das 32 entregas concluídas, recuperação da reserva expirada, descarte da confirmação antiga e 256 entregas finais contadas. [CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42905965564). Alteração somente de teste, sem nova imagem do servidor.
- [x] [Cloud #42](https://github.com/joaosouz4dev/wppconnect-cloud/pull/42): o guia público explica como renovar o QR na mesma sessão, deduplicar por X-WPP-Event-ID e encaminhar eventos a workflows. Passaram 50 testes, build, TypeScript, design e navegador nos dois idiomas. O [guia em produção](https://cloud.wppconnect.io/docs) retornou 200 e os três títulos localizados em inglês e português na revisão 72a54cec957de711459e2073a4e98ca28649675f. [CI da principal](https://github.com/joaosouz4dev/wppconnect-cloud/actions/runs/42907505395).
- [x] Chatwoot, somente piloto local: quatro testes Go com HTTP e SQLite verificaram uma resposta pública de atendente para um contato previamente vinculado, assinatura e escopo, concorrência, restart, deduplicação e recuperação por consulta sem repetir o envio. Sem credenciais ou destinatário real, sem implantação. O código do piloto permanece local; a ação conjunta de publicação foi bloqueada pela revisão automática.
- [ ] Chatwoot completo: provisionamento automático de contatos e conversas, anexos nos dois sentidos e aceite ponta a ponta com o telefone. O marco de texto abaixo registra a evolução do piloto.
- [ ] Mantidos os adiamentos do usuário: nenhum monitoramento ou canal de alertas por ora; a segunda conta e os testes em Android/iPhone/Web ficam para depois.

## Chatwoot e texto entre uma sessão e uma inbox de API

- [x] [Servidor Go #29](https://github.com/wppconnect-team/wppconnect-server-go/pull/29): adaptador publicado na principal. A imagem oficial do Chatwoot 4.19.1 assinou um callback com a chave padrão de 24 caracteres do Rails. O teste cobriu confirmação, repetição e alteração do corpo. [Execução com a imagem oficial](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42910842345).
- [x] [Servidor Go #30](https://github.com/wppconnect-team/wppconnect-server-go/pull/30): a aplicação Chatwoot 4.19.1, com PostgreSQL e Redis, criou conta, contato, conversa e inbox temporários. O adaptador processou o callback de atendente assinado a partir da mensagem da aplicação. No outro sentido, recebeu duas entregas sintéticas de texto e a API confirmou uma mensagem persistida. [Teste da inbox e do assinador](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42911977000).
- [x] Casos de falha com HTTP e SQLite locais: conta, inbox ou contato divergentes, notas privadas, ecos, mídia, texto temporário, perda de confirmação, deduplicação por ID da mensagem, recuperação após reinício e concorrência. O caso de confirmação perdida reconcilia por consulta ao destino. O teste com a inbox real usa um serviço simulado para a operação de envio no WhatsApp.

- O registro do piloto local descreve a primeira tentativa. A ferramenta rejeitou um comando composto com a mensagem genérica `blocked by policy`; a regra exata não foi informada. A branch e o PR #29 foram aceitos posteriormente pelo GitHub e o código foi incorporado à principal no SHA `217800f1b1ae6d2d1dc52d755d138ba0e8ccf94c`.

- [x] Publicação Go `c6dfe06d392f` e canary aprovados. Três reservas saudáveis, conexão existente preservada e containers temporários removidos em 2026-11-10T22:23:12 UTC. O processo do adaptador é configurado à parte pelo operador. [Recibo da pool e do teste Chatwoot](docs/chatwoot-text-acceptance.json).
- [x] [Cloud #43](https://github.com/joaosouz4dev/wppconnect-cloud/pull/43) publicada em `b2b9188987c9`. O guia de configuração passou na suíte da Cloud, no build e no navegador em inglês (1440 px) e português (375 px). Os dois títulos e o link de configuração também foram verificados no endereço público.

## Provisionamento do Chatwoot e escrita de catálogo

- [x] [Servidor Go #31](https://github.com/wppconnect-team/wppconnect-server-go/pull/31) incorporado. Chatwoot 4.19.1, PostgreSQL e Redis confirmaram que a perda da resposta de criação de contato, conversa ou vínculo de inbox é recuperável após reinício. O contato com o número exato foi reutilizado. [Execução do teste](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42915535317).
- [x] [Servidor Go #32](https://github.com/wppconnect-team/wppconnect-server-go/pull/32) implementa Z270, Z274 e Z287 com contratos nativos para o vendedor da sessão. Testes HTTP, SQLite e SDK simulado verificaram conflito de chave, isolamento da consulta de resultado, resposta incerta e visibilidade falsa explícita.
- [x] [Cloud #44](https://github.com/joaosouz4dev/wppconnect-cloud/pull/44) adiciona as três operações e a consulta de resultado sob `commerce:write`. O guia foi conferido em inglês e português no Chromium. A suíte local passou com 51 testes da Cloud e quatro de permissões do control plane.
- [ ] Aceite com vendedor de teste e produtos descartáveis para confirmar moderação e visibilidade após a mudança.

- [x] [CI do catálogo](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42917050226) com race, integrações e imagens amd64/arm64. A imagem de `75d34044d1ec` passou no canary e a pool recebeu três reservas Go saudáveis. A conexão existente ficou ativa. As três escritas foram recusadas na alocação de teste desconectada e seus resultados ficaram recuperáveis sem envio remoto.
- [x] A documentação da [Cloud em produção](https://cloud.wppconnect.io/docs) contém o novo guia em português e inglês. As quatro rotas exigiram autenticação, e uma rota desconhecida devolveu 404. [Recibo do catálogo, pool e integrações](docs/catalog-write-acceptance.json).
- [x] [Servidor Go #33](https://github.com/wppconnect-team/wppconnect-server-go/pull/33) incorporado. Com Redis real, 32 eventos de duas sessões esgotaram duas tentativas enquanto o webhook recebeu uma cópia de cada. Após reabrir o SQLite e restaurar o stream, o replay recuperou os 32 eventos respeitando sessão e destino. [Execução com os serviços](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42917781388/job/130184552313).

- [x] Recontagem dos contratos atuais de Node e Go e inventário dos sete diferenciais restantes do recorte em que o Node já oferecia a operação central. [Comparação, prioridades e critérios](docs/parity-remaining.md#comparação-atualizada-do-servidor-node-com-o-go).

- [x] [Servidor Go #34](https://github.com/wppconnect-team/wppconnect-server-go/pull/34) incorporado. Após a perda da confirmação do upload, o arquivamento explícito verifica os bytes, a geração do objeto e a reserva da sessão. O ensaio após reinício rejeitou conteúdo alterado ou corrompido, e o [teste com MinIO real](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42918653529/job/130187540165) confirmou um único upload.
- [ ] Arquivamento automático por sessão. Faltam a seleção de anexos persistentes no recebimento e a fila com recuperação após falha.

- [x] A revisão de recuperação S3 `43b76a436cea` passou no canary e chegou às três reservas Go. A auditoria de 23:06 UTC encontrou as três saudáveis e preservou a conexão anterior. Os objetos do teste ficaram no bucket temporário do MinIO no CI.
- [x] O verificador da matriz executa as quatro funções de teste HTTP que cobrem os 15 contratos mapeados e exige aprovação explícita de cada uma. Um ensaio negativo com uma fixture marcada como pulada foi reprovado pelo verificador.

- [x] [Servidor Go #35](https://github.com/wppconnect-team/wppconnect-server-go/pull/35) incorporado. Texto estendido marcado para visualização única perdeu a elegibilidade para encaminhamento ao Chatwoot. O histórico omite o marcador exclusivo do evento ao vivo, inclusive em linhas antigas. [CI do servidor](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42919986810) e [integração Chatwoot 4.19.1](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42919986661).

- [x] Duas inicializações isoladas da imagem Zuck v2.4.2 pelo digest AMD64 da matriz, cada uma com PostgreSQL descartável pronto, sem portas públicas, volumes de clientes ou telefones. Ambas encerraram com código 1 antes dos testes HTTP. O diagnóstico da aplicação informou `LICENSE_KEY` não configurada. Os dois containers de cada tentativa foram removidos. [Recibo da referência](docs/reference-runtime-acceptance.json).
- [ ] Comparação comportamental das APIs lado a lado. Requer uma instalação de homologação da referência já licenciada ou uma licença válida para a instância de teste. A ausência da chave não interfere na execução dos testes nativos do Go, da Cloud, do Chatwoot, do Redis e do MinIO.

- [x] Pool final em `72c35ee4d7ce`, três reservas Go saudáveis em 2026-11-10T23:17:07.675798+00:00. O primeiro canary terminou com erro antes da promoção. A repetição isolada passou com a mesma imagem. O binário leu uma linha de histórico de versão antiga, preservou o texto e omitiu o marcador exclusivo de mensagens ao vivo. A sessão pareada permaneceu conectada. [Recibo da publicação](docs/catalog-write-acceptance.json).

## Consulta real ao catálogo da sessão pareada

- [x] Consultar catálogo e coleções na sessão autorizada. As duas leituras retornaram HTTP 502. O SDK identificou GraphQL 400 nas duas chamadas e não obteve um perfil comercial identificável para a conta. O total de produtos permanece desconhecido. [Evidência de leitura](docs/catalog-live-account-acceptance.json).
- [x] A correção do [servidor Go #36](https://github.com/wppconnect-team/wppconnect-server-go/pull/36) passou no CI e foi verificada pelo Manager na sessão pareada. Com a revisão `f739440a12d7`, as duas consultas retornaram `WhatsApp catalog query failed (GraphQL 400)`. A sessão permaneceu conectada e as três reservas Go continuaram saudáveis na revisão `72c35ee4d7ce`. [CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42936768438).
- [ ] Confirmar leitura de um produto descartável, edição de coleção e visibilidade em uma conta WhatsApp Business de teste. O operador confirmou que a sessão usada é WhatsApp comum e não tem catálogo. Nenhuma escrita foi enviada a essa sessão.
- [ ] Comparar respostas HTTP com a Zuck v2.4.2. O operador confirmou que não há instância licenciada nem chave de teste disponível.

## Carrinho do vendedor (Z266)

- [x] Implementar `POST /api/{session}/set-cart-enabled` no Go. O SDK verifica se o WhatsApp devolveu o booleano solicitado. A rota exige chave de idempotência e o vendedor vem da sessão. Os testes cobriram os dois valores, resposta incerta, repetição, conflito e consulta de resultado em outra sessão. [Servidor Go #37](https://github.com/wppconnect-team/wppconnect-server-go/pull/37) e [CI](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42938509260).
- [x] Atualizar a matriz reproduzível e executar as quatro funções de teste HTTP associadas aos 16 contratos nativos registrados sobre a referência de 303 operações. O Go descreve 162 operações e o Node descreve 158. A interseção por método e caminho passou a 82. [Comparação de contratos](docs/parity-remaining.md#comparação-atualizada-do-servidor-node-com-o-go).
- [x] Publicar a imagem Go `bffcd90760a3` para novas instâncias e conferir as três reservas saudáveis. A instância descartável exigiu chave e recuperou o resultado das quatro escritas comerciais pelo Manager. A sessão pessoal ficou conectada. [Recibo da publicação](docs/merchant-cart-acceptance.json).
- [x] Publicar o escopo `commerce:write` e o guia da Cloud para o carrinho. A suíte passou com 51 testes e o navegador conferiu o conteúdo em inglês e português. Em produção, as duas traduções responderam 200 e as rotas de escrita e consulta de resultado exigiram autenticação. [Cloud #45](https://github.com/joaosouz4dev/wppconnect-cloud/pull/45) e [guia público](https://cloud.wppconnect.io/docs).
- [ ] Confirmar a configuração do carrinho em um vendedor de teste e verificar o resultado como comprador. A conta autorizada para mensagens usa o WhatsApp comum e não possui catálogo.

## Socket do WhatsApp Web e favoritos (Z136)

- [x] Medir endpoints e contagem de frames em Chrome sem pareamento. Duas conexões a `/ws/chat`, nas portas 443 e 5222. A primeira trocou dados na amostra. [Estudo e reprodução](docs/web-socket.md) e [metadados observados](docs/web-socket-observation.json).
- [x] Validar a remontagem de frames com cabeçalho dividido e reduzir de 73.728 para zero os bytes alocados por operação no benchmark de 64 KiB. [HyperMeow #14](https://github.com/wppconnect-team/hypermeow/pull/14).
- [x] Ampliar a automação diária de 09:30 de São Paulo para comparar o cliente público, o transporte do SDK e a revisão consumida pelo servidor, com fixtures sem conta real. A primeira observação está publicada neste plano.
- [x] Ativar os dois endpoints no server-go para sessões novas e restauradas, preservando o cliente HTTP da sessão. Vinte repetições locais verificaram que uma recusa do upgrade de qualquer endpoint permite receber um frame no outro. [Servidor Go #38](https://github.com/wppconnect-team/wppconnect-server-go/pull/38) e [CI do servidor](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42948935834).
- [x] Z136 na conta pessoal autorizada. Uma mensagem de teste foi favoritada, a sessão persistiu o evento observado, reconectou sem QR depois de reiniciar e confirmou a desfavoritação. O estado inicial foi restaurado. O verificador da matriz executou cinco funções de teste para 17 contratos nativos. [Recibo do aceite](docs/socket-star-acceptance.json).
- [x] Publicar o servidor na pool, conferir as três reservas Go saudáveis e remover as alocações descartáveis. A Cloud publicou os escopos de leitura e escrita de favoritos. Os 52 testes da Cloud passaram, o guia abriu em português e inglês, e as três operações recusaram a chamada pública sem credencial. [Cloud #46](https://github.com/joaosouz4dev/wppconnect-cloud/pull/46) e [resultado do teste](docs/socket-star-acceptance.json).

## Cinco diferenciais do Node com aceite isolado

- [x] Z290 - Remover canal com confirmação de identidade, papel de dono e proteção contra repetição.
- [x] Z271 - Excluir coleções com confirmação agregada e limite de lote.
- [x] Z282 - Criar produto preservando preço, moeda, mídia e resposta de revisão.
- [x] Z285 - Editar produto e conferir o ID devolvido.
- [x] Z283 - Excluir produtos e representar a contagem parcial.
- [x] Conferir escopos da Cloud, contratos OpenAPI e resultados persistidos nas suítes isoladas.
- [ ] Aceite externo adiado. A exclusão de canal foi limitada a testes isolados pelo operador. O catálogo exige conta Business e a comparação HTTP exige Zuck licenciada.

O verificador executou sete funções de teste HTTP para 22 contratos nativos mapeados sobre 303 operações de referência. Os dez diferenciais inicialmente apontados na comparação Node/Go têm testes de contrato. O CI do SDK executou as suítes de Go em duas versões. Na Cloud, os 52 testes locais passaram, assim como TypeScript, build e a leitura da documentação no navegador nas larguras 320, 375, 768 e 1440. A verificação desta fase usou exclusivamente recursos isolados. [SDK #16](https://github.com/wppconnect-team/hypermeow/pull/16), [servidor #39](https://github.com/wppconnect-team/wppconnect-server-go/pull/39), [Cloud #47](https://github.com/joaosouz4dev/wppconnect-cloud/pull/47) e [comparação dos cinco contratos](docs/parity-remaining.md#cinco-contratos-em-relação-ao-swagger-da-referência).

- [x] Encerrar e reiniciar um processo Go descartável com os mesmos arquivos SQLite. Os cinco comandos da sessão desconectada recuperaram os recibos originais após o reinício, e todos recusaram uma mudança de corpo com a chave anterior. [Teste reproduzível](docs/check-isolated-mutation-replay.py).

- [x] Incorporar as mudanças na principal do HyperMeow (`ba182bad1b6a`), do servidor Go (`1293a52d5bfd`) e da Cloud (`58a5af83a4c7`). Passaram a detecção de corridas em Go, os testes das imagens AMD64/ARM64 e a validação da Cloud. O ensaio não usou a sessão pareada nem promoveu uma imagem para a pool. [CI do servidor](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42957781408) e [CI da Cloud](https://github.com/joaosouz4dev/wppconnect-cloud/actions/runs/42957749748).

- [x] Conferir a imagem publicada `wppconnect/wppconnect-server-go:v0.1.35`. O manifesto contém `linux/amd64` e `linux/arm64` e tem digest `sha256:d1079e6ab33d1c20896d67d4ddbef6de0a77ef30242a73c209d97e1baca160b5`. [Release v0.1.35](https://github.com/wppconnect-team/wppconnect-server-go/releases/tag/v0.1.35) e [execução da publicação](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42958534363).

- [x] Atualizar a pool da Cloud para a v0.1.35. Em 11/11/2026 às 01h43 de São Paulo, as três reservas Go estavam saudáveis no digest publicado. Uma nova alocação pelo Manager confirmou as verificações de cinco comandos sem pareamento. As quatro alocações de teste foram removidas e as duas instâncias já alocadas preservaram o mesmo contêiner e o mesmo horário de início. A configuração e a imagem anteriores ficaram disponíveis para reversão. [Registro da atualização](docs/cloud-pool-v0.1.35.json).
