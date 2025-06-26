#include <bpfcore/vmlinux.h>
#include <bpfcore/bpf_helpers.h>
#include <common/map_sizing.h>
#include <logger/bpf_dbg.h>
// #include "pid.h"
#include <bpfcore/bpf_tracing.h>

#include <generictracer/k_send_receive.h>

// #include "sockaddr.h"
// #include "tcp_info.h"
// #include "k_tracer_defs.h"
// #include "http_ssl_defs.h"
// #include "pin_internal.h"
// #include "k_send_receive.h"
// #include "k_unix_sock.h"

char __license[] SEC("license") = "Dual MIT/GPL";

typedef struct capability_info {
    int cap;
    u64 pid;
} capability_info_t;

const capability_info_t *unused_2 __attribute__((unused));

// Temporary tracking of capabilities
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 12);
} capability_events SEC(".maps");

SEC("kprobe/capable")
int BPF_KPROBE(beyla_kprobe_capable, int cap) {
    u64 pid = bpf_get_current_pid_tgid();

    capability_info_t *trace = bpf_ringbuf_reserve(&capability_events, sizeof(capability_info_t), 0);
    if (trace) {
        //TODO: Log the system time. Can we use bpf_ktime_get_tai_ns?
        // https://docs.ebpf.io/linux/helper-function/bpf_ktime_get_tai_ns/
        // bpf_map_update_elem(&capability_events, &id, &cap, BPF_ANY);

        trace->cap = cap;
        trace->pid = pid;
        bpf_dbg_printk("=== capable updating ring buffer ===");
        bpf_ringbuf_submit(trace, 0);
    }

    return 0;
}
