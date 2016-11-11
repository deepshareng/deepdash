function addDate(dadd, gran) {
    var a = new Date();
    a = a.valueOf();
    var b;
    var month;
    switch (gran) {
        case "d":
            a = a + dadd * 24 * 60 * 60 * 1000;
            a = new Date(a);
            month = a.getMonth() + 1;
            b = a.getDate() + "/" + month;
            break;
        case "w":
            a = a + dadd * 24 * 60 * 60 * 1000 * 7;
            a = new Date(a);
            month = a.getMonth() + 1;
            b = a.getDate() + "/" + month;
            break;
        case "m":
            a = new Date(a);
            month = (a.getMonth() + dadd + 12) % 12 + 1;
            b = month + "月";
            break;
    }
    return b;
}

function randomInt(max) {
    return 0;
    /*
    max = max || 1000;
    return Math.floor(Math.random() * max);
    */
}

function escapeHtml(text) {
    return text
        .replace(/&/g, '&amp;')
        .replace(/ /g, '&nbsp;')
        .replace(/\"/g, '&quot;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
}

function isFullUrl(url) {
    if (url && /^http/.test(url)) {
        return true;
    } else {
        return false;
    }
}

module.exports = {
    addDate: addDate,
    randomInt: randomInt,
    escapeHtml: escapeHtml,
    isFullUrl: isFullUrl,
};
