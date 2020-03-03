# dbexport
Pacote para exportar e sincronizar de forma fácil os objetos do seu banco de dados.

Adicione na seção scripts do seu `package.json`

```
"scripts": {
    "db:export": "dbexport",
    "db:sync": "dbsync",
    ...
  },
```

#### rodando
`npm run db:export`

`npm run db:export procedures`

`npm run db:export triggers tg_pessoa_ins_after`

`npm run db:sync`
