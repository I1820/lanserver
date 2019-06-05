/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 04-06-2019
 * |
 * | File Name:     devices_collection_index.js
 * +===============================================
 */
/* eslint-env mongo */

var collection = "devices";

db[collection].createIndex({
  deveui: 1,
}, {
  unique: true,
});
