cd $1

if [ -f hosts ]; then
  cat hosts >> /etc/hosts
fi

exec ./bin/server