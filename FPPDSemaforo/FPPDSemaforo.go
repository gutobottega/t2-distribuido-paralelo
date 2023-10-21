// Disciplina de Modelos de Computacao Concorrente
// Escola Politecnica - PUCRS
// Prof.  Fernando Dotti
// ATENCAO - Instrucoes rapidas:
// COLOQUE O PACOTE FPPDSemaforo (este arquivo) DENTRO DE UM DIRETORIO(PASTA)
// CHAMADO FPPDSemaforo, NO DIRETORIO CORRENTE (onde esta seu codigo).
//
// No seu codigo que usa semaforo, faca:
// import (
//	"./FPPDSemaforo"
// )
// exemplo de declaracao de um semaforo:
//      FPPDSemaforo.Semaphore s = FPPDSemaforo.NewSemaphore(1)
//      s.Wait()
//      s.Signal()

package FPPDSemaforo

// ---------------------------------
type Semaphore struct { // este semáforo implementa quaquer numero de creditos em "v"
	v    int           // valor do semaforo: negativo significa proc bloqueado
	fila chan struct{} // canal para bloquear os processos se v < 0
	sc   chan struct{} // canal para atomicidade das operacoes wait e signal
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v:    init,                   // valor inicial de creditos
		fila: make(chan struct{}),    // canal sincrono para bloquear processos
		sc:   make(chan struct{}, 1), // usaremos este como semaforo para SC, somente 0 ou 1
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sc <- struct{}{} // SC do semaforo feita com canal
	s.v--              // decrementa valor
	if s.v < 0 {       // se negativo era 0 ou menor, tem que bloquear
		<-s.sc               // antes de bloq, libera acesso
		s.fila <- struct{}{} // bloqueia proc
	} else {
		<-s.sc // libera acesso
	}
}

func (s *Semaphore) Signal() {
	s.sc <- struct{}{} // entra sc
	s.v++
	if s.v <= 0 { // tem processo bloqueado ?
		<-s.fila // desbloqueia
	}
	<-s.sc // libera SC para outra op
}

// ----------------------------------
