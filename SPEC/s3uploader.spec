# TODO: need to test

%define         debug_package %{nil}

Name:           clickhouse-proxy-auth
Version:        1.0.0
Release:        1%{?dist}.vteam
Summary:        Access verification service for CH clusters

Group:          Applications/Internet
License:        unknown
Source0:        https://github.com/13excite/clickhouse-proxy-auth/archive/v%{version}.tar.gz


BuildArch:      x86_64
BuildRequires:  golang >= 1.18


%description

%prep

%build
make build

%install

%{__mkdir} -p $RPM_BUILD_ROOT/usr/local/bin
%{__mkdir} -p $RPM_BUILD_ROOT%{_sysconfdir}
%{__mkdir} -p $RPM_BUILD_ROOT%{_unitdir}
%{__install} -m 744 -p ./clickhouse-proxy-auth \
    $RPM_BUILD_ROOT/usr/local/bin/clickhouse-proxy-auth
%{__install} -m 644 -p ./config.yaml.example \
    $RPM_BUILD_ROOT%{_sysconfdir}/ch-proxy-auth.yaml
%{__install} -m644 ./unit.service \
    $RPM_BUILD_ROOT%{_unitdir}/clickhouse-proxy-auth.service

%files
%defattr(-,www-data,www-data,-)
%attr(0744, www-data, www-data) /usr/local/bin/clickhouse-proxy-auth
%attr(644, www-data, www-data) %{_sysconfdir}/ch-proxy-auth.yaml
%config(noreplace) %{_sysconfdir}/ch-proxy-auth.yaml


%doc

%changelog
*  Sun Oct 1 2023 Vladimir Dulenov
- initial rpm release
