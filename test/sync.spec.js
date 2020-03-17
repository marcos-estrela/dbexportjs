expect = require('chai').expect

const DbSync = require('../src/sync')

describe('DbSync', () => {
  describe('sync for views', () => {
    it('test getDropQueryFromFileName', async () => {
      const sync = new DbSync('mysql')
      let content = await sync.getDropQueryFromFileName('export/views/vw_simples.sql')
      expect(content).to.equal('DROP view IF EXISTS vw_simples;')
    })

    // it('test addDropQueryIfNotExists', async () => {
    //   const sync = new DbSync('mysql')
    //   const sql = 'CREATE OR REPLACE VIEW vw_simples AS SELECT id, name FROM simplestb'
    //   content = await sync.addDropQueryIfNotExists(sql)
    //   expect(content).to.equal('DROP view IF EXISTS vw_simples;');
    // })

    it('test getDropQuery', async () => {
      const sync = new DbSync('mysql')
      const sql = 'CREATE OR REPLACE VIEW vw_simples AS SELECT id, name FROM simplestb;'
      content = await sync.getDropQuery(sql)
      expect(content).to.equal('');
    })
  })

  describe('sync for procedures', () => {
    it('test getDropQueryFromFileName', async () => {
      const sync = new DbSync('mysql')
      let content = await sync.getDropQueryFromFileName('export/procedures/sp_simples_ins.sql')
      expect(content).to.equal('DROP procedure IF EXISTS sp_simples_ins;')
    })

    // it('test addDropQueryIfNotExists', async () => {
    //   const sync = new DbSync('mysql')
    //   const sql = 'CREATE OR REPLACE VIEW vw_simples AS SELECT id, name FROM simplestb'
    //   content = await sync.addDropQueryIfNotExists(sql)
    //   expect(content).to.equal('DROP view IF EXISTS vw_simples;');
    // })

    it('test getDropQuery', async () => {
      const sync = new DbSync('mysql')
      const sql = `CREATE PROCEDURE sp_ola()
                  COMMENT 'Procedure para executar olá'
                  begin
                      select 'olá george william barros de';
                  end`
      content = await sync.getDropQuery(sql)
      expect(content).to.equal('DROP PROCEDURE IF EXISTS sp_ola;\n\n');
    })
  })
})
