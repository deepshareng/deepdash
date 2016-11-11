var MarketingReducer = function(state={}, action) {
    if (Object.keys(state).length === 0 && state.constructor === Object) {
        return {
            metrics: [],
            tableHeads: [],
            tableBodyKeys: [],
            tableData: [],
            dateStartStr: '',
            dateEndStr: '',
        };
    } else {
        if (action.type === 'UPDATE_MARKETING_INFO') {
            return Object.assign({}, state, action.conf);
        } else {
            return state;
        }
    }
};

module.exports = MarketingReducer;
