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
	print(cmd)
	os.execute(cmd)
end

gitpush()
