import random
import subprocess
import threading

import pytest


def _write_random(f, seed, length):
    rng = random.Random()
    rng.seed(seed)
    f.write(rng.randbytes(length))


def _read_random(f, seed, length):
    rng = random.Random()
    rng.seed(seed)
    remaining = length
    while remaining > 0:
        buff = f.read(remaining)
        if len(buff) == 0:
            break
        remaining -= len(buff)
        assert buff == rng.randbytes(len(buff))
    assert remaining == 0


def alice(process: subprocess.Popen[str]):
    _write_random(process.stdin, 1, 5000)
    process.stdin.close()
    process.wait()


def bob(process: subprocess.Popen[str]):
    process.stdin.close()
    _read_random(process.stdout, 1, 5000)
    process.wait()


# TODO
def make_arg_tuples(template):
    remote_tuples = [("tcp-client,localhost:8000", "tcp-server,:8000")]
    result = []
    for alice_args, bob_args in template:
        for alice_remote, bob_remote in remote_tuples:
            pass  # TODO: replace


@pytest.mark.parametrize("alice_args, bob_args", [
    (["--send", "tcp-server,:8002"], ["--receive", "tcp-client,localhost:8002"]),
    (["--send", "tcp-client,localhost:8003"], ["--receive", "tcp-server,:8003"]),
    (["tcp-server,:8000"], ["tcp-client,localhost:8000"]),
    (["tcp-client,localhost:8001"], ["tcp-server,:8001"]),
])
def test_cat(alice_args, bob_args):
    alice_process = subprocess.Popen(
        ["tunme", "cat"] + alice_args,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=2,
    )
    alice_thread = threading.Thread(target=alice, args=(alice_process,))
    alice_thread.start()

    bob_process = subprocess.Popen(
        ["tunme", "cat"] + bob_args,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=2,
    )
    bob_thread = threading.Thread(target=bob, args=(bob_process,))
    bob_thread.start()

    alice_thread.join(5)
    bob_thread.join(5)
    assert not alice_thread.is_alive()
    assert not bob_thread.is_alive()
