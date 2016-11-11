var simulatorConf = {
    'IOS9_YYB_INSTALLED': {
        platform: 'IOS9',
        download: 'YYB',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/f541f6f21c034b6a53a2efdc344d2f44b8d265c7/embed',
    },
    'IOS9_YYB_NEW': {
        platform: 'IOS9',
        download: 'YYB',
        install: 'NEW',
        url: 'https://modao.cc/app/9e3671ea2501d52e80c26b2794abd37869f7a003/embed',
    },
    'IOS9_DIRECT_INSTALLED': {
        platform: 'IOS9',
        download: 'DIRECT',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/BlityVaQCll63Z81pwkxO8d4mkFM7Xp/embed',
    },
    'IOS9_DIRECT_NEW': {
        platform: 'IOS9',
        download: 'DIRECT',
        install: 'NEW',
        url: 'https://modao.cc/app/2ae2a7e270b0291be2b41675dde004a80bf10e31/embed',
    },

    'IOS8_YYB_INSTALLED': {
        platform: 'IOS8',
        download: 'YYB',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/de125e893597b4fa545be1b0ac7de186471d2caf/embed',
    },
    'IOS8_YYB_NEW': {
        platform: 'IOS8',
        download: 'YYB',
        install: 'NEW',
        url: 'https://modao.cc/app/54c7d631a7e7c11341c7f8426e8414701906cc30/embed',
    },

    'IOS8_DIRECT_INSTALLED': {
        platform: 'IOS8',
        download: 'DIRECT',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/4a21154d69ba2a24d74eebd8e0887083f298aaba/embed',
    },
    'IOS8_DIRECT_NEW': {
        platform: 'IOS8',
        download: 'DIRECT',
        install: 'NEW',
        url: 'https://modao.cc/app/34739a9620de811885b30acf570ace2db995507b/embed',
    },

    'ANDROID_YYB_INSTALLED': {
        platform: 'ANDROID',
        download: 'YYB',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/6a5464b2d7697d9a9ec95deec0ff6b4fe097c1e9/embed',
    },
    'ANDROID_YYB_NEW': {
        platform: 'ANDROID',
        download: 'YYB',
        install: 'NEW',
        url: 'https://modao.cc/app/b111b906e4e50be4f6d5ca84e19e7f3ffa6fc9ec/embed',
    },
    'ANDROID_DIRECT_INSTALLED': {
        platform: 'ANDROID',
        download: 'DIRECT',
        install: 'INSTALLED',
        url: 'https://modao.cc/app/b6b0c6d4ad3f5c0c3494e26ac473a435c2742573/embed',
    },
    'ANDROID_DIRECT_NEW': {
        platform: 'ANDROID',
        download: 'DIRECT',
        install: 'NEW',
        url: 'https://modao.cc/app/b2ae92d704662888c3b41eb24f77050257e7dbd8/embed',
    },
};

var SimulatorReducer = function(state={}, action) {
    // the logic of update components!!
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return simulatorConf['IOS9_YYB_INSTALLED'];
    } else {
        if (action.type === 'UPDATE_SIMULATOR_TYPE') {
            var simulatorType = (action.platform || state.platform) + '_' + (action.download || state.download) + '_' + (action.install || state.install);
            return Object.assign({}, state, simulatorConf[simulatorType]);
        } else {
            return state;
        }
    }
};

module.exports = {
    simulatorConf: simulatorConf,
    SimulatorReducer: SimulatorReducer,
};
