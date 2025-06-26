// Do we need a "go:build linux"? Beyla is Linux-only anyway.

package capabilitytracer

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log/slog"

	"github.com/cilium/ebpf"

	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/app/request"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/beyla"
	beyla_ebpf "github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/ebpf"
	ebpfcommon "github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/ebpf/common"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/ebpf/ringbuf"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/exec"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/goexec"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/components/svc"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/config"
	"github.com/open-telemetry/opentelemetry-ebpf-instrumentation/pkg/pipe/msg"
)

//go:generate $BPF2GO -cc $BPF_CLANG -cflags $BPF_CFLAGS -type capability_info_t -target amd64,arm64 Bpf ../../../../bpf/capabilitytracer/capability_tracer.c -- -I../../../../bpf
//go:generate $BPF2GO -cc $BPF_CLANG -cflags $BPF_CFLAGS -type capability_info_t -target amd64,arm64 BpfDebug ../../../../bpf/capabilitytracer/capability_tracer.c -- -I../../../../bpf -DBPF_DEBUG

type BPFCapabilityInfo BpfCapabilityInfoT

type Tracer struct {
	pidsFilter ebpfcommon.ServiceFilter
	cfg        *beyla.Config
	bpfObjects BpfObjects
	closers    []io.Closer
	log        *slog.Logger
}

// AddInstrumentedLibRef implements ebpf.Tracer.
func (p *Tracer) AddInstrumentedLibRef(uint64) {
}

func (p *Tracer) AllowPID(pid, ns uint32, svc *svc.Attrs) {
}

func (p *Tracer) BlockPID(pid, ns uint32) {
}

// AlreadyInstrumentedLib implements ebpf.Tracer.
func (p *Tracer) AlreadyInstrumentedLib(uint64) bool {
	return false
}

// Constants implements ebpf.Tracer.
func (p *Tracer) Constants() map[string]any {
	return nil
}

// GoProbes implements ebpf.Tracer.
func (p *Tracer) GoProbes() map[string][]*ebpfcommon.ProbeDesc {
	return nil
}

// ProcessBinary implements ebpf.Tracer.
func (p *Tracer) ProcessBinary(*exec.FileInfo) {
}

// RecordInstrumentedLib implements ebpf.Tracer.
func (p *Tracer) RecordInstrumentedLib(uint64, []io.Closer) {
}

// RegisterOffsets implements ebpf.Tracer.
func (p *Tracer) RegisterOffsets(*exec.FileInfo, *goexec.Offsets) {
}

// SockMsgs implements ebpf.Tracer.
func (p *Tracer) SockMsgs() []ebpfcommon.SockMsg {
	return nil
}

// SockOps implements ebpf.Tracer.
func (p *Tracer) SockOps() []ebpfcommon.SockOps {
	return nil
}

// SocketFilters implements ebpf.Tracer.
func (p *Tracer) SocketFilters() []*ebpf.Program {
	return nil
}

// UProbes implements ebpf.Tracer.
func (p *Tracer) UProbes() map[string]map[string][]*ebpfcommon.ProbeDesc {
	return nil
}

// UnlinkInstrumentedLib implements ebpf.Tracer.
func (p *Tracer) UnlinkInstrumentedLib(uint64) {
}

var _ beyla_ebpf.Tracer = (*Tracer)(nil)

func New(cfg *beyla.Config) *Tracer {
	log := slog.With("component", "capabilitytracer.Tracer")
	return &Tracer{
		log:        log,
		cfg:        cfg,
		pidsFilter: ebpfcommon.CommonPIDsFilter(&cfg.Discovery),
	}
}

func (p *Tracer) Load() (*ebpf.CollectionSpec, error) {
	loader := LoadBpf
	if p.cfg.EBPF.BpfDebug {
		loader = LoadBpfDebug
	}

	return loader()
}

func (p *Tracer) BpfObjects() any {
	return &p.bpfObjects
}

func (p *Tracer) AddCloser(c ...io.Closer) {
	p.closers = append(p.closers, c...)
}

func (p *Tracer) KProbes() map[string]ebpfcommon.ProbeDesc {
	kprobes := map[string]ebpfcommon.ProbeDesc{
		"capable": {
			Required: true,
			Start:    p.bpfObjects.BeylaKprobeCapable,
		},
	}

	return kprobes
}

func (p *Tracer) Tracepoints() map[string]ebpfcommon.ProbeDesc {
	return nil
}

func (p *Tracer) SetupTailCalls() {}

func (p *Tracer) Run(ctx context.Context, ebpfEventContext *ebpfcommon.EBPFEventContext, eventsChan *msg.Queue[[]request.Span]) {
	ebpfcommon.ForwardRingbuf(
		&p.cfg.EBPF,
		p.bpfObjects.CapabilityEvents,
		&ebpfcommon.IdentityPidsFilter{},
		p.process,
		p.log,
		nil,
		append(p.closers, &p.bpfObjects)...,
	)(ctx, eventsChan)
}

func (p *Tracer) process(_ *ebpfcommon.EBPFParseContext, _ *config.EBPFTracer, record *ringbuf.Record, _ ebpfcommon.ServiceFilter) (request.Span, bool, error) {
	var event BPFCapabilityInfo

	p.log.Info("capabilitytracer::process start")

	err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event)
	if err != nil {
		p.log.Info("capabilitytracer::process failed to parse")
		return request.Span{}, true, err
	}

	p.log.Info("capabilitytracer::process parsed capability", "cap", event.Cap)

	return request.Span{}, false, nil
}

func (p *Tracer) Required() bool {
	return false
}
