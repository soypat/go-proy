n = 100;
D = 3;
X = rand(D,n);
w = rand(D,1);
y_NoNoise = X.'*w;


mean = 0;%1/(1+1/999)^999; %1/e
sigma_n = 2;
noise =mean+ sigma_n.*randn(n,1);

y = y_NoNoise+noise;

%% Calc. Covariance
% xMeans = sum(X,2)/n;
% COV = nan(D,D);
% COV2 = nan(D,D);
% for i = 1:D
%    for j = 1:D
%        COV(i,j) = sum( (X(i,:)-xMeans(i)) .* (X(j,:)-xMeans(j)) )/(n-1);
%        COV2(i,j) = sum( (X(i,:)-xMeans(i)) .* (X(j,:)-xMeans(j)) )/(n-1);
% 
%    end
% end
COV = getCovariance(X);
% COV %Diag [s=1 -> d=0.8] [s=.1 -> d= ]
% COV = exp(-.5*(xp-xq).^2);

w_bar = (sigma_n^(-2))*( (sigma_n^(-2)) * X*X.' + inv(COV))^-1*X*y;
w_bar2 =(sigma_n^(-2))*( (sigma_n^(-2)) * X*X.' )^-1*X*y;

LinFunc = getLinFun(D);

f = X.'*w_bar;
f2 = X.'*w_bar2;
err= abs(f-y).^2;
err2 = abs(f2-y).^2;
plot(err)
hold on
plot(err2)



function [covariance] = getCovariance(X)
    % X (D,n)
    [~,n] = size(X);
    X=X.';
%     Iv = ones(n,1);
    DeviationX = X - sum(X)/n;
    covariance = DeviationX.'*DeviationX/n;
    return
end

function [linearFunction] = getLinFun(Dimension)
functionString = "@(w,x) ";

for i = 1:Dimension
    functionString = strcat(functionString,"w(",string(i),")*x(",string(i),")" );
    if i~=Dimension
        functionString = strcat(functionString,"+");
    end
end
linearFunction = eval(functionString);
end

% i_p =randi([1 n]);
% i_q = randi([1 n]);
% while i_p == i_q
%     i_q = randi([1 n]);
% end
% xp = X(:,i_p);
% xq = X(:,i_q);