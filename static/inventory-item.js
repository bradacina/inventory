'use strict'

Vue.component('inventory-item',{
    props:['item'],

    template: `<div>
    <div>{{item.Title}}</div>
    <div>{{item.SKU}}</div>
    <div>{{item.Barcode}}</div>
    <div>{{item.Quantity}}</div>
    </div>`
});