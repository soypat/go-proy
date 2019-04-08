%% Toy 1 dimension input problem

% f = @(x) w(1)+w(2)*x+w(3)*x^2;
deg = 2;
n=5;

D=1;
%% DATOS
f = @(x) 4+2.*x+x.^2; %No se conoce, pero la uso para general el vector y
sigma_n = 2;
noise = sigma_n*randn(D,n); %Ruido en el problema

% xtest = linspace(-1,6,n);
xtest = 2.5; %#ok<NBRAK>
xstar = xtest*ones(n,1);

X = (linspace(0,5,n)).';%4*(rand(D,n)+1);
xvec = linspace(0,5,40);
fvec = feval(f,xvec);
phi= @(x) [1,x,x.^2].'; %mi modelo/funcionalidad. (ojo el transpuesto)
N=length(phi(rand(1,1)));

w = rand(N,1);

Phi = zeros(N,n);
y = zeros(n,1);

for j = 1:n
    Phi(:,j) = phi(X(j));%D>1 => phi(X(1,j))
    y(j) = f(X(j))+noise(j);
end
scatter(X,y)
K=getKernel(X,X);
L=chol(K+sigma_n^2*eye(size(K)));
alpha = L.'\(L\y);
k_star = getKernel(X,xstar);
fbar = L\k_star;


%Otra forma de calcular fbar
fbar2 = k_star*((K+sigma_n^2*eye(size(K)))\y);
scatter(xstar,fbar2,'*');
hold on
% scatter(xtest,fbar,'.')
hold on
scatter(X,y,'o')
hold on
plot(xvec,fvec)
function [kernel] = getKernel(Xp,Xq)
[~,np]=size(Xp);
[~,nq]=size(Xq);
if nq~=np
    kernel = nan(np);
    return
end
kernel = zeros(np);
for i = 1:np
   for j=1:nq
       xp=Xp(:,i);
       xq=Xq(:,j);
       kernel(i,j) = exp(-.5*norm(xp-xq)^2);
   end
end
end

function [pol] = getPolynomialFunction(deg)
    str="@(x,w) ";
    for i=0:deg
        str = str+"w("+sprintf("%i",i+1)+")*x.^("+sprintf("%i",i)+")";
        if i==deg
            str=str+";";
            continue
        end
            str=str+"+";
    end
    pol = eval(str);
end
