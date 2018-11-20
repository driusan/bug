package bugapp

type ArgumentList []string

func (args ArgumentList) HasArgument(arg string) bool {
	for _, argCandidate := range args {
		if arg == argCandidate {
			return true
		}
	}
	return false
}

func (args ArgumentList) GetArgument(argname, defaultVal string) string {
	for idx, argCandidate := range args {
		if argname == argCandidate {
			// If it's the last argument, then return string
			// "true" because we can't return idx+1, but it
			// shouldn't be the default value when the argument
			// isn't provided either..
			if idx >= len(args)-1 {
				return "true"
			}
			return args[idx+1]
		}
	}
	return defaultVal
}

func (args ArgumentList) GetAndRemoveArguments(argnames []string) (ArgumentList, []string) {
	var nextArgumentType int = -1
	matches := make([]string, len(argnames))
	var retArgs []string
	for idx, argCandidate := range args {
		// The last token was in argnames, so this one
		// is the value. Set it in matches and reset
		// nextArgumentType and continue to the next
		// possible token
		if nextArgumentType != -1 {
			matches[nextArgumentType] = argCandidate
			nextArgumentType = -1
			continue
		}

		// Check if this is a argname we're looking for
		for argidx, argname := range argnames {
			if argname == argCandidate {
				if idx >= len(args)-1 {
					matches[argidx] = "true"
				}
				nextArgumentType = argidx
				break
			}
		}

		// It wasn't an argname, so add it to the return
		if nextArgumentType == -1 {
			retArgs = append(retArgs, argCandidate)
		}
	}
	return retArgs, matches
}
