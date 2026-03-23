function encode(strs) {
    return strs.map(s => `${s.length}#${s}`).join("");
}

function decode(s) {
    const result = [];
    let i = 0;
    while (i < s.length) {
        let j = s.indexOf("#", i);
        const length = parseInt(s.substring(i, j), 10);
        result.push(s.substring(j + 1, j + 1 + length));
        i = j + 1 + length;
    }
    return result;
}

module.exports = { encode, decode };
