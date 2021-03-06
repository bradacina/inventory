'use strict'

Vue.component('inventory-items-list', {
    template: `
    <div>
        <inventory-item v-for="item in items" :key="item.SKU" v-bind:item="item" />
    </div>`,

    data: function() {
        return {items: window.inventoryItems};
    }
});