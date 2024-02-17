import json

with open('./docs/opcodes.json', 'r') as fp:
    spec = json.load(fp)


# Print unprefixed cycles

ops_unprefixed = list(spec['unprefixed'].values())

cycles = [min(op['cycles']) // 4 for op in ops_unprefixed]

assert len(cycles) == 256

print("== Unprefixed M-Cycle Times == ")
for i in range(0, len(cycles), 16):
    print(f'/* {i:0X} */'.replace('0', 'x'), ', '.join(map(str, cycles[i:i+16])), ', //')
print("==")

_bytes = [op['bytes'] for op in ops_unprefixed]
print("== Unprefixed Op Lens == ")
for i in range(0, len(_bytes), 16):
    print(f'/* {i:0X} */'.replace('0', 'x'), ', '.join(map(str, _bytes[i:i+16])), ', //')
print("==")

def get_name(op):
    operands = ''
    if op['operands']:
        operands = ':'+','.join(o['name'] for o in op['operands'])
    return f"\"{op['mnemonic']}{operands}\"" 

names = [get_name(op) for op in ops_unprefixed]

for i in range(0, len(cycles), 16):
    print(', '.join(map(str, names[i:i+16])))