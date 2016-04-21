local cs = "hello lua"
local dir = "/root/Doc/bin/huaban"
-- local dir = "warms"

function makecommand( )
	local one = "cd " .. dir
	local two = [[git add .;git commit -m "init";git push]]
	return one .. ";" .. two
end

function gitpush( )
	local cmd = makecommand()
	print(os.date("%Y-%m-%d %H:%M:%S"), " Push")
	local r = os.execute(cmd)
	if r == true then
	    print(os.date("%Y-%m-%d %H:%M:%S"), " : Push OK")
	else
	    print(os.date("%Y-%m-%d %H:%M:%S"), " : Push Faile")
	end
end

gitpush()
