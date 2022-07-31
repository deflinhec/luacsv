local DATA = DATA --csv 數據
local logger = require 'logger' -- 日誌模塊
local json = require 'json' -- JSON模塊

local type = type
local print = print
local pairs = pairs
local assert = assert

assert(type(DATA)=="table", "should be table")
assert(type(logger.debug)=="function", "should have function")
assert(type(logger.info)=="function", "should have function")
assert(type(logger.warn)=="function", "should have function")
assert(type(logger.error)=="function", "should have function")
assert(type(json.encode)=="function", "should have function")
assert(type(json.decode)=="function", "should have function")

for _, v in pairs(DATA) do
    assert(type(v)=="table", "should be table")
    logger.info("info-message", v)
    print(json.encode(v))
end