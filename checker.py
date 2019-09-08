import json

check_dict = dict()


class node:
    def __init__(self, data):
        self.trace_id = data["id"]
        self.start = data["root"]["start"]
        self.end = data["root"]["end"]
        self.service = data["root"]["service"]
        self.span = data["root"]["span"]
        self.calls = data["root"]["calls"]

    def compare(self, other):
        assert other.start == self.start
        assert other.end == self.end
        assert other.service == self.service
        assert other.span == self.span
        assert len(other.calls) == len(self.calls), f"Expected: {len(other.calls)}, got {len(self.calls)}. Id: {self.trace_id}"


expected = open(...)
actual = open(...)

for line in expected.readlines():
    tree = json.loads(line)
    check_dict[tree["id"]] = node(tree)

for line in actual.readlines():
    tree = json.loads(line)

    check_dict[tree["id"]].compare(node(tree))
    del check_dict[tree["id"]]

assert len(check_dict) == 0, "Ensure that all trace ids were found"
