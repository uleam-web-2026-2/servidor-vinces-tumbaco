package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/uleam-web-2026-2/KEMA/internal/avisos"
	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"github.com/uleam-web-2026-2/KEMA/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type registroAviso struct {
	usuario string
	aviso   avisos.Aviso
}

type emisorPrueba struct {
	enviados []registroAviso
}

func (e *emisorPrueba) Enviar(usuario string, aviso avisos.Aviso) {
	e.enviados = append(e.enviados, registroAviso{
		usuario: usuario,
		aviso:   aviso,
	})
}

type estadoSimulado struct {
	mu           sync.Mutex
	id           int32
	estado       string
	ofreceID     sql.NullInt32
	recibeID     sql.NullInt32
	fallarUpdate bool
}

var (
	simMu     sync.Mutex
	simActual *estadoSimulado
)

func init() {
	sql.Register("simulado_kema", &driverSimulado{})
}

type driverSimulado struct{}

func (d *driverSimulado) Open(name string) (driver.Conn, error) {
	return &conexionSimulada{}, nil
}

type conexionSimulada struct{}

func (c *conexionSimulada) Prepare(query string) (driver.Stmt, error) {
	return &sentenciaSimulada{query: query}, nil
}

func (c *conexionSimulada) Close() error { return nil }
func (c *conexionSimulada) Begin() (driver.Tx, error) {
	return &transaccionSimulada{}, nil
}

type transaccionSimulada struct{}

func (t *transaccionSimulada) Commit() error   { return nil }
func (t *transaccionSimulada) Rollback() error { return nil }

type sentenciaSimulada struct {
	query string
}

func (s *sentenciaSimulada) Close() error  { return nil }
func (s *sentenciaSimulada) NumInput() int { return -1 }

func (s *sentenciaSimulada) Exec(args []driver.Value) (driver.Result, error) {
	return s.ExecContext(context.Background(), toNamedValues(args))
}

func (s *sentenciaSimulada) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(context.Background(), toNamedValues(args))
}

func toNamedValues(args []driver.Value) []driver.NamedValue {
	out := make([]driver.NamedValue, len(args))
	for i, v := range args {
		out[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}
	return out
}

func (s *sentenciaSimulada) ExecContext(_ context.Context, _ []driver.NamedValue) (driver.Result, error) {
	simMu.Lock()
	sim := simActual
	simMu.Unlock()
	if sim == nil {
		return nil, errors.New("sin estado simulado")
	}

	sim.mu.Lock()
	defer sim.mu.Unlock()

	if sim.fallarUpdate {
		return nil, errors.New("fallo simulado al guardar estado")
	}

	q := strings.ToUpper(s.query)
	if strings.Contains(q, "UPDATE") && strings.Contains(q, "INTERCAMBIOS") {
		sim.estado = "aceptado"
		return driver.RowsAffected(1), nil
	}

	return driver.RowsAffected(1), nil
}

func (s *sentenciaSimulada) QueryContext(_ context.Context, _ []driver.NamedValue) (driver.Rows, error) {
	simMu.Lock()
	sim := simActual
	simMu.Unlock()
	if sim == nil {
		return nil, errors.New("sin estado simulado")
	}

	sim.mu.Lock()
	defer sim.mu.Unlock()

	return &filasIntercambio{
		entregado: false,
		id:        sim.id,
		estado:    sim.estado,
		ofreceID:  sim.ofreceID,
		recibeID:  sim.recibeID,
	}, nil
}

type filasIntercambio struct {
	entregado bool
	id        int32
	estado    string
	ofreceID  sql.NullInt32
	recibeID  sql.NullInt32
}

func (f *filasIntercambio) Columns() []string {
	return []string{"id", "estado", "usuario_ofrece_id", "usuario_recibe_id"}
}

func (f *filasIntercambio) Close() error { return nil }

func (f *filasIntercambio) Next(dest []driver.Value) error {
	if f.entregado {
		return io.EOF
	}
	f.entregado = true
	dest[0] = int64(f.id)
	dest[1] = f.estado
	if f.ofreceID.Valid {
		dest[2] = int64(f.ofreceID.Int32)
	} else {
		dest[2] = nil
	}
	if f.recibeID.Valid {
		dest[3] = int64(f.recibeID.Int32)
	} else {
		dest[3] = nil
	}
	return nil
}

func prepararServicioSimulado(t *testing.T, emisor EmisorAvisos, sim *estadoSimulado) *IntercambiosService {
	t.Helper()

	simMu.Lock()
	simActual = sim
	simMu.Unlock()

	sqlDB, err := sql.Open("simulado_kema", "")
	if err != nil {
		t.Fatalf("no se pudo abrir sql simulado: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir gorm con conexion simulada: %v", err)
	}

	repo := repository.NuevoIntercambiosRepository(db, sqlDB)
	return NuevoIntercambiosService(repo, emisor)
}

func TestCambiarEstadoEnviaAvisoSoloAParticipantes(t *testing.T) {
	emisor := &emisorPrueba{}
	sim := &estadoSimulado{
		id:       1,
		estado:   intercambios.EstadoPendiente,
		ofreceID: sql.NullInt32{Int32: 10, Valid: true},
		recibeID: sql.NullInt32{Int32: 20, Valid: true},
	}
	srv := prepararServicioSimulado(t, emisor, sim)

	actualizado, err := srv.CambiarEstado(1, intercambios.EstadoAceptado)
	if err != nil {
		t.Fatalf("cambio de estado inesperado: %v", err)
	}
	if actualizado.Estado != intercambios.EstadoAceptado {
		t.Fatalf("estado obtenido %s, esperado %s", actualizado.Estado, intercambios.EstadoAceptado)
	}

	if len(emisor.enviados) != 2 {
		t.Fatalf("se esperaban 2 avisos, se enviaron %d", len(emisor.enviados))
	}
	if emisor.enviados[0].usuario != "10" || emisor.enviados[1].usuario != "20" {
		t.Fatalf("destinatarios inesperados: %+v", emisor.enviados)
	}
	for _, env := range emisor.enviados {
		if env.aviso.ID != 1 || env.aviso.Estado != intercambios.EstadoAceptado {
			t.Fatalf("aviso inesperado: %+v", env.aviso)
		}
	}
}

func TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia(t *testing.T) {
	emisor := &emisorPrueba{}
	sim := &estadoSimulado{
		id:           1,
		estado:       intercambios.EstadoPendiente,
		ofreceID:     sql.NullInt32{Int32: 10, Valid: true},
		recibeID:     sql.NullInt32{Int32: 20, Valid: true},
		fallarUpdate: true,
	}
	srv := prepararServicioSimulado(t, emisor, sim)

	_, err := srv.CambiarEstado(1, intercambios.EstadoAceptado)
	if err == nil {
		t.Fatal("se esperaba error al guardar el estado")
	}
	if len(emisor.enviados) != 0 {
		t.Fatalf("no debian enviarse avisos tras error de guardado, llegaron %d", len(emisor.enviados))
	}
}

func TestCambiarEstadoConParticipantesNulosNoNotifica(t *testing.T) {
	emisor := &emisorPrueba{}
	sim := &estadoSimulado{
		id:       1,
		estado:   intercambios.EstadoPendiente,
		ofreceID: sql.NullInt32{Valid: false},
		recibeID: sql.NullInt32{Valid: false},
	}
	srv := prepararServicioSimulado(t, emisor, sim)

	if _, err := srv.CambiarEstado(1, intercambios.EstadoAceptado); err != nil {
		t.Fatalf("cambio de estado inesperado: %v", err)
	}
	if len(emisor.enviados) != 0 {
		t.Fatalf("no debian enviarse avisos con participantes nulos, llegaron %d", len(emisor.enviados))
	}
}
